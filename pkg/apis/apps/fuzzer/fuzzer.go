package fuzzer

import (
	"github.com/LiangNing7/minerx/pkg/apis/apps"

	fuzz "github.com/google/gofuzz"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the apps api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(j *apps.Chain, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again

			// match defaulting
		},
		func(j *apps.MinerSet, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again

			// match defaulting
			if len(j.Spec.Selector.MatchLabels) == 0 && len(j.Spec.Selector.MatchExpressions) == 0 {
				j.Spec.Selector = metav1.LabelSelector{MatchLabels: j.Spec.Template.Labels}
			}
			if len(j.Labels) == 0 {
				j.Labels = j.Spec.Template.Labels
			}
		},
		func(j *apps.Miner, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again

			// match defaulting
		},
	}
}
