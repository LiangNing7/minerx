package options

import (
	"net"

	genericapiserver "k8s.io/apiserver/pkg/server"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kube-openapi/pkg/common"

	controlplaneoptions "github.com/LiangNing7/minerx/internal/controlplane/apiserver/options"
	"github.com/LiangNing7/minerx/pkg/apiserver/storage"
)

const defaultEtcdPathPrefix = "/registry/liangning7.cn"

type ServerRunOptions struct {
	*controlplaneoptions.Options

	Extra
}

type Extra struct {
	MasterCount int
	// In the future, perhaps an "minerlet" will be added, similar to the "kubelet".
	// MinerxletConfig minerletclient.MinerxletClientConfig.
	APIServerServiceIP     net.IP
	EndpointReconclierType string

	// For external resources.
	ExternalRESTStorageProviders []storage.RESTStorageProvider
	ExternalPostStartHooks       map[string]genericapiserver.PostStartHookFunc
	GetOpenAPIDefinitions        common.GetOpenAPIDefinitions
}

// NewServerRunOptions returns a new ServerRunOptions.
func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{
		Options: controlplaneoptions.NewOptions(),
		Extra: Extra{
			MasterCount:            1,
			ExternalPostStartHooks: make(map[string]genericapiserver.PostStartHookFunc),
		},
	}

	return o
}

func (o ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	o.Options.AddFlags(&fss)

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")

	fs.IntVar(&o.MasterCount, "apiserver-count", o.MasterCount,
		"The number of apiservers running in the cluster, must be a positive number. (In use when --endpoint-reconciler-type=master-count is enabled.)")
	fs.MarkDeprecated("apiserver-count", "apiserver-count is deprecated and will be removed in a future version.")

	return fss
}
