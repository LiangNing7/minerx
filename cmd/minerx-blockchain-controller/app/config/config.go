package config

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/LiangNing7/minerx/internal/controller/blockchain/apis/config"
)

// Config is the main context object for the controller.
type Config struct {
	ComponentConfig *config.BlockchainControllerConfiguration

	// the rest config for the master
	Kubeconfig *restclient.Config

	// Kubernetes clientset used to create miner pods.
	ProviderClient kubernetes.Interface

	ProviderCluster cluster.Cluster
}
