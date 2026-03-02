package latest

import (
	"github.com/LiangNing7/minerx/internal/controller/job/apis/config"
	"github.com/LiangNing7/minerx/internal/controller/job/apis/config/scheme"
	"github.com/LiangNing7/minerx/internal/controller/job/apis/config/v1beta1"
)

// Default creates a default configuration of the latest versioned type.
// This function needs to be updated whenever we bump the job controller's component config version.
func Default() (*config.JobControllerConfiguration, error) {
	versionedCfg := v1beta1.JobControllerConfiguration{}

	scheme.Scheme.Default(&versionedCfg)
	cfg := config.JobControllerConfiguration{}
	if err := scheme.Scheme.Convert(&versionedCfg, &cfg, nil); err != nil {
		return nil, err
	}
	// We don't set this field in internal/controller/job/apis/config/{version}/conversion.go
	// because the field will be cleared later by API machinery during
	// conversion. See JobControllerConfiguration internal type definition for
	// more details.
	cfg.TypeMeta.APIVersion = v1beta1.SchemeGroupVersion.String()
	return &cfg, nil
}
