package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	componentbasevalidation "k8s.io/component-base/config/validation"

	"github.com/LiangNing7/minerx/internal/controller/apis/config"
	genericvalidation "github.com/LiangNing7/minerx/pkg/config/validation"
)

// Validate ensures validation of the MinerControllerConfiguration struct.
func Validate(cc *config.MinerXControllerManagerConfiguration) field.ErrorList {
	allErrs := field.ErrorList{}

	newPath := field.NewPath("MinerXControllerManagerConfiguration")

	allErrs = append(allErrs, componentbasevalidation.ValidateLeaderElectionConfiguration(&cc.Generic.LeaderElection, newPath.Child("generic", "leaderElection"))...)
	allErrs = append(allErrs, genericvalidation.ValidateMySQLConfiguration(&cc.MySQL, newPath.Child("mysql"))...)
	allErrs = append(allErrs, genericvalidation.ValidateGenericControllerManagerConfiguration(&cc.Generic, newPath.Child("generic"))...)

	return allErrs
}
