package ensurer

import (
	flowcontrolv1 "k8s.io/api/flowcontrol/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	flowcontrolapisv1 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1"

	flowcontrolclient "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned/typed/flowcontrol/v1"
	flowcontrollisters "github.com/LiangNing7/minerx/pkg/generated/listers/flowcontrol/v1"
)

func NewPriorityLevelConfigurationOps(client flowcontrolclient.PriorityLevelConfigurationInterface, lister flowcontrollisters.PriorityLevelConfigurationLister) ObjectOps[*flowcontrolv1.PriorityLevelConfiguration] {
	return NewObjectOps[*flowcontrolv1.PriorityLevelConfiguration](client, lister, (*flowcontrolv1.PriorityLevelConfiguration).DeepCopy,
		plcReplaceSpec, plcSpecEqualish)
}

func plcReplaceSpec(into, from *flowcontrolv1.PriorityLevelConfiguration) *flowcontrolv1.PriorityLevelConfiguration {
	copy := into.DeepCopy()
	copy.Spec = *from.Spec.DeepCopy()
	return copy
}

func plcSpecEqualish(expected, actual *flowcontrolv1.PriorityLevelConfiguration) bool {
	copiedExpected := expected.DeepCopy()
	flowcontrolapisv1.SetObjectDefaults_PriorityLevelConfiguration(copiedExpected)
	if expected.Name == flowcontrolv1.PriorityLevelConfigurationNameExempt {
		if actual.Spec.Exempt == nil {
			return false
		}
		copiedExpected.Spec.Exempt.NominalConcurrencyShares = actual.Spec.Exempt.NominalConcurrencyShares
		copiedExpected.Spec.Exempt.LendablePercent = actual.Spec.Exempt.LendablePercent
	}
	return equality.Semantic.DeepEqual(copiedExpected.Spec, actual.Spec)
}
