//nolint:gocritic
package options

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	"github.com/LiangNing7/minerx/internal/controller/apis/config"
	"github.com/LiangNing7/minerx/internal/controller/apis/config/scheme"
	"github.com/LiangNing7/minerx/internal/controller/apis/config/v1beta1"
)

func LoadConfigFromFile(file string) (*config.MinerXControllerManagerConfiguration, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return loadConfig(data)
}

func loadConfig(data []byte) (*config.MinerXControllerManagerConfiguration, error) {
	// The UniversalDecoder runs defaulting and returns the internal type by default.
	obj, gvk, err := scheme.Codecs.UniversalDecoder().Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}
	if cfgObj, ok := obj.(*config.MinerXControllerManagerConfiguration); ok {
		// We don't set this field in pkg/scheduler/apis/config/{version}/conversion.go
		// because the field will be cleared later by API machinery during
		// conversion. See MinerXControllerManagerConfiguration internal type definition for
		// more details.
		cfgObj.TypeMeta.APIVersion = gvk.GroupVersion().String()
		return cfgObj, nil
	}
	return nil, fmt.Errorf("couldn't decode as MinerXControllerManagerConfiguration, got %s: ", gvk)
}

// LogOrWriteConfig logs the completed component config and writes it into the given file name as YAML, if either is enabled.
func LogOrWriteConfig(fileName string, cfg *config.MinerXControllerManagerConfiguration) error {
	if !(klog.V(2).Enabled() || len(fileName) > 0) {
		return nil
	}

	buf, err := encodeConfig(cfg)
	if err != nil {
		return err
	}

	if klog.V(2).Enabled() {
		klog.Info("Using component config", "config", buf.String())
	}

	if len(fileName) > 0 {
		configFile, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer configFile.Close()
		if _, err := io.Copy(configFile, buf); err != nil {
			return err
		}
		klog.InfoS("Wrote configuration", "file", fileName)
		os.Exit(0)
	}
	return nil
}

func encodeConfig(cfg *config.MinerXControllerManagerConfiguration) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return buf, fmt.Errorf("unable to locate encoder -- %q is not a supported media type", mediaType)
	}

	var encoder runtime.Encoder
	switch cfg.TypeMeta.APIVersion {
	case v1beta1.SchemeGroupVersion.String():
		encoder = scheme.Codecs.EncoderForVersion(info.Serializer, v1beta1.SchemeGroupVersion)
	default:
		encoder = scheme.Codecs.EncoderForVersion(info.Serializer, v1beta1.SchemeGroupVersion)
	}
	if err := encoder.Encode(cfg, buf); err != nil {
		return buf, err
	}
	return buf, nil
}
