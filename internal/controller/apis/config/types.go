package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genericconfig "github.com/LiangNing7/minerx/pkg/config"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MinerXControllerManagerConfiguration contains elements describing minerx-controller manager.
type MinerXControllerManagerConfiguration struct {
	metav1.TypeMeta

	// FeatureGates is a map of feature names to bools that enable or disable alpha/experimental features.
	// FeatureGates map[string]bool

	// MySQL defines the configuration of mysql client.
	MySQL genericconfig.MySQLConfiguration `json:"mysql,omitempty"`

	// Generic holds configuration for a generic controller-manager
	Generic genericconfig.GenericControllerManagerConfiguration

	// GarbageCollectorControllerConfiguration holds configuration for
	// GarbageCollectorController related features.
	GarbageCollectorController genericconfig.GarbageCollectorControllerConfiguration

	// ChainControllerConfiguration holds configuration for ChainController related features.
	ChainController ChainControllerConfiguration
}

type ChainControllerConfiguration struct {
	// Image specify the blockchain node image.
	Image string
}
