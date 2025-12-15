package serializer

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// workaround serializer to solve bug https://github.com/kubernetes/kubernetes/issues/86666
type protocolShieldSerializers struct {
	*serializer.CodecFactory
	accepts []runtime.SerializerInfo
}

// NewProtocolShieldSerializers deletion fails because object does not implement protobuf marshaling interface.
// protocolShieldSerializers is a workaround serializer to solve this problem.
// Reference: https://github.com/kubernetes/kubernetes/issues/86666.
func NewProtocolShieldSerializers(codecs *serializer.CodecFactory) *protocolShieldSerializers {
	if codecs == nil {
		return nil
	}
	pss := &protocolShieldSerializers{
		CodecFactory: codecs,
		accepts:      []runtime.SerializerInfo{},
	}
	pss.accepts = append(pss.accepts, codecs.SupportedMediaTypes()...)
	return pss
}

// SupportedMediaTypes returns the RFC2046 media types that this factory has serializers for.
func (pss *protocolShieldSerializers) SupportedMediaTypes() []runtime.SerializerInfo {
	return pss.accepts
}

// EncoderForVersion returns an encoder that targets the provided group version.
func (pss *protocolShieldSerializers) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return pss.CodecFactory.EncoderForVersion(encoder, gv)
}

// DecoderToVersion returns a decoder that targets the provided group version.
func (pss *protocolShieldSerializers) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return pss.CodecFactory.DecoderToVersion(decoder, gv)
}
