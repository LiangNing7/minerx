// Package options contains flags and options for initializing an apiserver
package options

import (
	"strings"

	"github.com/LiangNing7/goutils/pkg/app"
	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/LiangNing7/minerx/internal/pkg/feature"
	"github.com/LiangNing7/minerx/internal/pump"
)

const (
	// UserAgent is the userAgent name when starting onex-pump server.
	UserAgent = "minerx-pump"
)

// ServerOptions contains the configuration options for the server.
type ServerOptions struct {
	HealthOptions *genericoptions.HealthOptions `json:"health" mapstructure:"health"`
	KafkaOptions  *genericoptions.KafkaOptions  `json:"kafka" mapstructure:"kafka"`
	MongoOptions  *genericoptions.MongoOptions  `json:"mongo" mapstructure:"mongo"`
	FeatureGates  map[string]bool               `json:"feature-gates"`
}

// Ensure ServerOptions implements the app.NamedFlagSetOptions interface.
var _ app.NamedFlagSetOptions = (*ServerOptions)(nil)

// NewServerOptions creates a ServerOptions instance with default values.
func NewServerOptions() *ServerOptions {
	o := &ServerOptions{
		HealthOptions: genericoptions.NewHealthOptions(),
		KafkaOptions:  genericoptions.NewKafkaOptions(),
		MongoOptions:  genericoptions.NewMongoOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *ServerOptions) Flags() (fss cliflag.NamedFlagSets) {
	o.HealthOptions.AddFlags(fss.FlagSet("health"))
	o.KafkaOptions.AddFlags(fss.FlagSet("kafka"))
	o.MongoOptions.AddFlags(fss.FlagSet("mongo"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *ServerOptions) Complete() error {
	url := o.MongoOptions.URL
	if !strings.HasPrefix(url, "mongodb://") && !strings.HasPrefix(url, "mongodb+srv://") {
		// Preserve backwards compatibility for hostnames without a
		// scheme, broken in go 1.8. Remove in Telegraf 2.0
		o.MongoOptions.URL = "mongodb://" + o.MongoOptions.URL
	}

	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate checks whether the options in ServerOptions are valid.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	errs = append(errs, o.HealthOptions.Validate()...)
	errs = append(errs, o.KafkaOptions.Validate()...)
	errs = append(errs, o.MongoOptions.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// Config builds an pump.Config based on ServerOptions.
func (o *ServerOptions) Config() (*pump.Config, error) {
	return &pump.Config{
		KafkaOptions: o.KafkaOptions,
		MongoOptions: o.MongoOptions,
	}, nil
}
