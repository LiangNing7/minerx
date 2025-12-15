package ensurer

import (
	flowcontrolv1 "k8s.io/api/flowcontrol/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	flowcontrolapisv1 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1"

	flowcontrolclient "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned/typed/flowcontrol/v1"
	flowcontrollisters "github.com/LiangNing7/minerx/pkg/generated/listers/flowcontrol/v1"
)

func NewFlowSchemaOps(client flowcontrolclient.FlowSchemaInterface, cache flowcontrollisters.FlowSchemaLister) ObjectOps[*flowcontrolv1.FlowSchema] {
	return NewObjectOps[*flowcontrolv1.FlowSchema](client, cache, (*flowcontrolv1.FlowSchema).DeepCopy, flowSchemaReplaceSpec, flowSchemaSpecEqual)
}

func flowSchemaReplaceSpec(into, from *flowcontrolv1.FlowSchema) *flowcontrolv1.FlowSchema {
	copy := into.DeepCopy()
	copy.Spec = *from.Spec.DeepCopy()
	return copy
}

func flowSchemaSpecEqual(expected, actual *flowcontrolv1.FlowSchema) bool {
	copiedExpectedSpec := expected.Spec.DeepCopy()
	flowcontrolapisv1.SetDefaults_FlowSchemaSpec(copiedExpectedSpec)
	return equality.Semantic.DeepEqual(copiedExpectedSpec, &actual.Spec)
}
