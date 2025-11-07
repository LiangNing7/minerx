package install

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"

	batchfuzzer "github.com/LiangNing7/minerx/pkg/apis/batch/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, batchfuzzer.Funcs)
	// TODO: enable protobuf generation for the sample-apiserver
	// roundtrip.RoundTripProtobufTestForAPIGroup(t, Install, wardlefuzzer.Funcs)
}
