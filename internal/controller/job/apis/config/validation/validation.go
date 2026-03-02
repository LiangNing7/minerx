package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	componentbasevalidation "k8s.io/component-base/config/validation"

	"github.com/LiangNing7/minerx/internal/controller/job/apis/config"
	genericvalidation "github.com/LiangNing7/minerx/pkg/config/validation"
)

// Validate ensures validation of the JobControllerConfiguration struct.
func Validate(cc *config.JobControllerConfiguration) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, componentbasevalidation.ValidateLeaderElectionConfiguration(&cc.Generic.LeaderElection, field.NewPath("generic", "leaderElection"))...)
	allErrs = append(allErrs, genericvalidation.ValidateGenericControllerManagerConfiguration(&cc.Generic, field.NewPath("generic"))...)

	return allErrs
}
