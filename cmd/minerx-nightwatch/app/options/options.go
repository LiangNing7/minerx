// Package options contains flags and options for initializing an nightwatch.
package options

import (
	"math"

	"github.com/LiangNing7/goutils/pkg/app"
	"github.com/LiangNing7/goutils/pkg/log"
	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	"github.com/LiangNing7/goutils/pkg/watch"
	"github.com/spf13/viper"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/LiangNing7/minerx/internal/nightwatch"
	"github.com/LiangNing7/minerx/internal/pkg/feature"
	kubeutil "github.com/LiangNing7/minerx/internal/pkg/util/kube"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
)

const (
	// UserAgent is the userAgent name when starting minerx-nightwatch server.
	UserAgent = "minerx-nightwatch"
)

// ServerOptions contains the configuration options for the server.
type ServerOptions struct {
	HealthOptions         *genericoptions.HealthOptions  `json:"health" mapstructure:"health"`
	MySQLOptions          *genericoptions.MySQLOptions   `json:"mysql" mapstructure:"mysql"`
	RedisOptions          *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	WatchOptions          *watch.Options                 `json:"nightwatch" mapstructure:"nightwatch"`
	HTTPOptions           *genericoptions.HTTPOptions    `json:"http" mapstructure:"http"`
	TLSOptions            *genericoptions.TLSOptions     `json:"tls" mapstructure:"tls"`
	UserWatcherMaxWorkers int64                          `json:"user-watcher-max-workers" mapstructure:"user-watcher-max-workers"`
	DisableRESTServer     bool                           `json:"disable-rest-server" mapstructure:"disable-rest-server"`
	Metrics               *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig   string          `json:"kubeconfig" mapstructure:"kubeconfig"`
	FeatureGates map[string]bool `json:"feature-gates"`
	Log          *log.Options    `json:"log" mapstructure:"log"`
}

// Ensure ServerOptions implements the app.NamedFlagSetOptions interface.
var _ app.NamedFlagSetOptions = (*ServerOptions)(nil)

// NewServerOptions creates a ServerOptions instance with default values.
func NewServerOptions() *ServerOptions {
	o := &ServerOptions{
		HealthOptions:         genericoptions.NewHealthOptions(),
		MySQLOptions:          genericoptions.NewMySQLOptions(),
		RedisOptions:          genericoptions.NewRedisOptions(),
		HTTPOptions:           genericoptions.NewHTTPOptions(),
		TLSOptions:            genericoptions.NewTLSOptions(),
		DisableRESTServer:     false,
		UserWatcherMaxWorkers: math.MaxInt64,
		WatchOptions:          watch.NewOptions(),
		Metrics:               genericoptions.NewMetricsOptions(),
		Log:                   log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *ServerOptions) Flags() (fss cliflag.NamedFlagSets) {
	o.HealthOptions.AddFlags(fss.FlagSet("health"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	o.WatchOptions.AddFlags(fss.FlagSet("watch"))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	o.Log.AddFlags(fss.FlagSet("log"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.BoolVar(&o.DisableRESTServer, "disable-rest-server", o.DisableRESTServer, "Disable the REST server functionality.")
	fs.Int64Var(&o.UserWatcherMaxWorkers, "user-watcher-max-workers", o.UserWatcherMaxWorkers, "Specify the maximum concurrency event of user watcher.")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *ServerOptions) Complete() error {
	if err := viper.Unmarshal(&o); err != nil {
		return err
	}

	if o.UserWatcherMaxWorkers < 1 {
		o.UserWatcherMaxWorkers = math.MaxInt64
	}

	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate checks whether the options in ServerOptions are valid.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	errs = append(errs, o.HealthOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.WatchOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// Config builds an nightwatch.Config based on ServerOptions.
func (o *ServerOptions) Config() (*nightwatch.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", o.Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeutil.SetDefaultClientOptions(kubeutil.AddUserAgent(kubeconfig, UserAgent))

	client, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	cfg := &nightwatch.Config{
		HealthOptions:         o.HealthOptions,
		MySQLOptions:          o.MySQLOptions,
		RedisOptions:          o.RedisOptions,
		HTTPOptions:           o.HTTPOptions,
		TLSOptions:            o.TLSOptions,
		WatchOptions:          o.WatchOptions,
		DisableRESTServer:     o.DisableRESTServer,
		UserWatcherMaxWorkers: o.UserWatcherMaxWorkers,
		Client:                client,
	}

	return cfg, nil
}
