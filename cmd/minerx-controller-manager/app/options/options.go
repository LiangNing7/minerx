// Package options provides the flags used for the controller manager.
package options

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/metrics"
	"k8s.io/kubernetes/pkg/controller/garbagecollector"

	controllermanagerconfig "github.com/LiangNing7/minerx/cmd/minerx-controller-manager/app/config"
	"github.com/LiangNing7/minerx/cmd/minerx-controller-manager/names"
	ctrlmgrconfig "github.com/LiangNing7/minerx/internal/controller/apis/config"
	"github.com/LiangNing7/minerx/internal/controller/apis/config/latest"
	clientcmdutil "github.com/LiangNing7/minerx/internal/pkg/util/clientcmd"
	kubeutil "github.com/LiangNing7/minerx/internal/pkg/util/kube"
	genericconfig "github.com/LiangNing7/minerx/pkg/config"
	genericconfigoptions "github.com/LiangNing7/minerx/pkg/config/options"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
)

const (
	// ControllerManagerUserAgent is the userAgent name when starting minerx-controller managers.
	ControllerManagerUserAgent = "minerx-controller-manager"
)

// Options is the main context object for the minerx-controller manager.
type Options struct {
	Generic                    *genericconfigoptions.GenericControllerManagerConfigurationOptions
	GarbageCollectorController *genericconfigoptions.GarbageCollectorControllerOptions
	MySQL                      *genericconfigoptions.MySQLOptions
	// NamespaceController *NamespaceControllerOptions

	// ConfigFile is the location of the miner controller server's configuration file.
	ConfigFile string

	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string

	// The address of the Kubernetes API server (overrides any value in kubeconfig).
	Master string
	// Path to kubeconfig with authoriaztion and master location information.
	Kubeconfig string
	Metrics    *metrics.Options
	Logs       *logs.Options

	// config is the minerx controller manager server's configuration object.
	// The default values.
	config *ctrlmgrconfig.MinerXControllerManagerConfiguration
}

// NewOptions creates a new Options with a default config.
func NewOptions() (*Options, error) {
	componentConfig, err := latest.Default()
	if err != nil {
		return nil, err
	}

	o := Options{
		Generic:                    genericconfigoptions.NewGenericControllerManagerConfigurationOptions(&componentConfig.Generic),
		GarbageCollectorController: genericconfigoptions.NewGarbageCollectorControllerOptions(&componentConfig.GarbageCollectorController),
		MySQL:                      genericconfigoptions.NewMySQLOptions(&componentConfig.MySQL),
		Kubeconfig:                 clientcmdutil.DefaultKubeconfig(),
		Metrics:                    metrics.NewOptions(),
		Logs:                       logs.NewOptions(),
	}

	gcIgnoredResources := make([]genericconfig.GroupResource, 0, len(garbagecollector.DefaultIgnoredResources()))
	for r := range garbagecollector.DefaultIgnoredResources() {
		gcIgnoredResources = append(gcIgnoredResources, genericconfig.GroupResource{Group: r.Group, Resource: r.Resource})
	}
	o.GarbageCollectorController.GCIgnoredResources = gcIgnoredResources
	o.Generic.LeaderElection.ResourceName = "minerx-controller-manager"
	o.Generic.LeaderElection.ResourceNamespace = metav1.NamespaceSystem

	return &o, nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags(allControllers []string, disabledControllers []string, controllerAliases map[string]string) cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	o.Generic.AddFlags(&fss, allControllers, disabledControllers, controllerAliases)
	o.GarbageCollectorController.AddFlags(fss.FlagSet(names.GarbageCollectorController))
	o.MySQL.AddFlags(fss.FlagSet("mysql"))

	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	logsapi.AddFlags(o.Logs, fss.FlagSet("logs"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	fs.StringVar(&o.WriteConfigTo, "write-config-to", o.WriteConfigTo, "If set, write the default configuration values to this file and exit.")
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")

	utilfeature.DefaultMutableFeatureGate.AddFlag(fss.FlagSet("generic"))

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {
	return nil
}

// ApplyTo fills up minerx controller manager config with options.
func (o *Options) ApplyTo(c *controllermanagerconfig.Config, allControllers []string, disabledControllers []string, controllerAliases map[string]string) error {
	if err := o.Generic.ApplyTo(&c.ComponentConfig.Generic, allControllers, disabledControllers, controllerAliases); err != nil {
		return err
	}
	if err := o.GarbageCollectorController.ApplyTo(&c.ComponentConfig.GarbageCollectorController); err != nil {
		return err
	}

	o.Metrics.Apply()

	return nil
}

// Validate is used to validate the options and config before launching the controller.
func (o *Options) Validate(allControllers []string, disabledControllers []string, controllerAliases map[string]string) error {
	var errs []error

	errs = append(errs, o.Generic.Validate(allControllers, disabledControllers, controllerAliases)...)
	errs = append(errs, o.GarbageCollectorController.Validate()...)

	// TODO: validate component config, master and kubeconfig

	return utilerrors.NewAggregate(errs)
}

// Config return a controller manager config objective.
func (o Options) Config(allControllers []string, disabledControllers []string, controllerAliases map[string]string) (*controllermanagerconfig.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(o.Master, o.Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeconfig.DisableCompression = true

	restConfig := kubeutil.AddUserAgent(kubeconfig, ControllerManagerUserAgent)
	client, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	c := &controllermanagerconfig.Config{
		Client:          client,
		Kubeconfig:      kubeutil.SetClientOptionsForController(restConfig),
		ComponentConfig: &ctrlmgrconfig.MinerXControllerManagerConfiguration{},
	}

	if err := o.ApplyTo(c, allControllers, disabledControllers, controllerAliases); err != nil {
		return nil, err
	}

	return c, nil
}
