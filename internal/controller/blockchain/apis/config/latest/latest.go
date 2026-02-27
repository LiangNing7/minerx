// Package latest provides version-agnostic access to the latest blockchain controller configuration.
package latest

import (
	"github.com/LiangNing7/minerx/internal/controller/blockchain/apis/config"
	"github.com/LiangNing7/minerx/internal/controller/blockchain/apis/config/scheme"
	"github.com/LiangNing7/minerx/internal/controller/blockchain/apis/config/v1beta1"
)

// Default creates a default configuration of the latest versioned type.
// This function needs to be updated whenever we bump the blockchain controller's component config version.
func Default() (*config.BlockchainControllerConfiguration, error) {
	versionedCfg := v1beta1.BlockchainControllerConfiguration{}

	scheme.Scheme.Default(&versionedCfg)
	cfg := config.BlockchainControllerConfiguration{}
	if err := scheme.Scheme.Convert(&versionedCfg, &cfg, nil); err != nil {
		return nil, err
	}
	// We don't set this field in internal/controller/blockchain/apis/config/{version}/conversion.go
	// because the field will be cleared later by API machinery during
	// conversion. See BlockchainControllerConfiguration internal type definition for
	// more details.
	cfg.APIVersion = v1beta1.SchemeGroupVersion.String()
	return &cfg, nil
}
