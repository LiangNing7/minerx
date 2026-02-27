package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"

	genericconfigv1beta1 "github.com/LiangNing7/minerx/pkg/config/v1beta1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_BlockchainControllerConfiguration sets additional defaults.
func SetDefaults_BlockchainControllerConfiguration(obj *BlockchainControllerConfiguration) {
	genericconfigv1beta1.RecommendedDefaultGenericControllerManagerConfiguration(&obj.Generic)

	if obj.FeatureGates == nil {
		obj.FeatureGates = make(map[string]bool)
	}

	if obj.ProviderKubeconfig == "" {
		// Here KUBECONFIG environment variable will not be used, KUBECONFIG is reserved for minerx-apiserver.
		obj.ProviderKubeconfig = clientcmd.RecommendedHomeFile
	}

	if len(obj.Types) == 0 {
		obj.Types = map[string]MinerProfile{
			"S1.SMALL1": {
				CPU:              resource.MustParse("50m"),
				Memory:           resource.MustParse("128Mi"),
				MiningDifficulty: 7,
			},
			"S1.SMALL2": {
				CPU:              resource.MustParse("100m"),
				Memory:           resource.MustParse("256Mi"),
				MiningDifficulty: 5,
			},
			"M1.MEDIUM1": {
				CPU:              resource.MustParse("150m"),
				Memory:           resource.MustParse("512Mi"),
				MiningDifficulty: 3,
			},
			"M1.MEDIUM2": {
				CPU:              resource.MustParse("250m"),
				Memory:           resource.MustParse("1024Mi"),
				MiningDifficulty: 1,
			},
		}
	}

	genericconfigv1beta1.RecommendedDefaultRedisConfiguration(&obj.Redis)
	genericconfigv1beta1.RecommendedDefaultMySQLConfiguration(&obj.MySQL)
	RecommendedDefaultChainControllerConfiguration(&obj.ChainController)
}

func RecommendedDefaultChainControllerConfiguration(obj *ChainControllerConfiguration) {
	if obj.Image == "" {
		obj.Image = "ccr.ccs.tencentyun.com/LiangNing7/minerx-toyblc-amd64:v0.1.0"
	}
}
