package latest

import (
	"github.com/LiangNing7/minerx/internal/controller/apis/config"
	"github.com/LiangNing7/minerx/internal/controller/apis/config/scheme"
	"github.com/LiangNing7/minerx/internal/controller/apis/config/v1beta1"
)

// Default creates a default configuration of the latest versioned type.
// This function needs to be updated whenever we bump the miner controller's component config version.
func Default() (*config.MinerXControllerManagerConfiguration, error) {
	versioned := v1beta1.MinerXControllerManagerConfiguration{}

	scheme.Scheme.Default(&versioned)
	cfg := config.MinerXControllerManagerConfiguration{}
	if err := scheme.Scheme.Convert(&versioned, &cfg, nil); err != nil {
		return nil, err
	}

	// We don't set this field in internal/controller/apis/config/{version}/conversion.go
	// because the field will be cleared later by API machinery during
	// conversion. See MinerControllerConfiguration internal type definition for
	// more details.
	cfg.TypeMeta.APIVersion = v1beta1.SchemeGroupVersion.String()
	return &cfg, nil
}
