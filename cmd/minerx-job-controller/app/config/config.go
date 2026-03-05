package config

import (
	restclient "k8s.io/client-go/rest"

	"github.com/LiangNing7/minerx/internal/controller/job/apis/config"
)

// Config is the main context object for the controller.
type Config struct {
	ComponentConfig *config.JobControllerConfiguration

	// the rest config for the master
	Kubeconfig *restclient.Config
}
