package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"

	genericconfigv1beta1 "github.com/LiangNing7/minerx/pkg/config/v1beta1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_MinerXControllerManagerConfiguration(obj *MinerXControllerManagerConfiguration) {
	genericconfigv1beta1.RecommendedDefaultGenericControllerManagerConfiguration(&obj.Generic)
	genericconfigv1beta1.RecommendedDefaultGarbageCollectorControllerConfiguration(&obj.GarbageCollectorController)
	RecommendedDefaultChainControllerConfiguration(&obj.ChainController)
}

func RecommendedDefaultChainControllerConfiguration(obj *ChainControllerConfiguration) {
	if obj.Image == "" {
		obj.Image = "ccr.ccs.tencentyun.com/LiangNing7/minerx-toyblc-amd64:v0.1.0"
	}
}
