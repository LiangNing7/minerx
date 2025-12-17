package initializer

import (
	"k8s.io/apiserver/pkg/admission"

	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
	"github.com/LiangNing7/minerx/pkg/generated/informers"
)

type pluginInitializer struct {
	informers informers.SharedInformerFactory
	client    clientset.Interface
	// authorizer        authorizer.Authorizer
	// featureGates      featuregate.FeatureGate
	stopCh <-chan struct{}
}

var _ admission.PluginInitializer = pluginInitializer{}

// New creates an instance of node admission plugins initializer.
func New(
	informers informers.SharedInformerFactory,
	client clientset.Interface,
) pluginInitializer {
	return pluginInitializer{
		informers: informers,
		client:    client,
	}
}

// Initialize checks the initialization interfaces implemented by a plugin
// and provide the appropriate initialization data.
func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsExternalInformerFactory); ok {
		wants.SetExternalInformerFactory(i.informers)
	}

	if wants, ok := plugin.(WantsExternalClientSet); ok {
		wants.SetExternalClientSet(i.client)
	}
}
