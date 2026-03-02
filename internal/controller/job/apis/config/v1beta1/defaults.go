package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"

	genericconfigv1beta1 "github.com/LiangNing7/minerx/pkg/config/v1beta1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_JobControllerConfiguration sets additional defaults.
func SetDefaults_JobControllerConfiguration(obj *JobControllerConfiguration) {
	genericconfigv1beta1.RecommendedDefaultGenericControllerManagerConfiguration(&obj.Generic)

	if obj.ConcurrentCronJobSyncs == 0 {
		obj.ConcurrentCronJobSyncs = 5
	}
	if obj.ConcurrentJobSyncs == 0 {
		obj.ConcurrentJobSyncs = 5
	}
}
