//nolint:unused
package gateway

import (
	"context"
	"os"

	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	"github.com/LiangNing7/goutils/pkg/server"
	"github.com/LiangNing7/goutils/pkg/version"
	"k8s.io/client-go/rest"

	"github.com/LiangNing7/minerx/internal/pkg/client/usercenter"
)

var (
	// Name is the name of the compiled software.
	Name = "minerx-gateway"

	ID, _ = os.Hostname()

	Version = version.Get().String()
)

// Config contains application-related configurations.
type Config struct {
	GRPCOptions       *genericoptions.GRPCOptions
	HTTPOptions       *genericoptions.HTTPOptions
	TLSOptions        *genericoptions.TLSOptions
	UserCenterOptions *usercenter.UserCenterOptions
	MySQLOptions      *genericoptions.MySQLOptions
	RedisOptions      *genericoptions.RedisOptions
	EtcdOptions       *genericoptions.EtcdOptions
	JaegerOptions     *genericoptions.JaegerOptions
	ConsulOptions     *genericoptions.ConsulOptions

	// the rest config for the onex-apiserver
	Kubeconfig *rest.Config
}

// Server represents the web server.
type Server struct {
	srv server.Server
}

// ServerConfig contains the core dependencies and configurations of the server.
type ServerConfig struct {
	cfg *Config
}

// NewServer initialized and returns a new Server instance.
func (cfg *Config) NewServer(ctx context.Context) (*Server, error) {
	_ = cfg
	return &Server{}, nil
}

// Run starts the server and listens for termination signals.
// It gracefully shuts down the server upon receiving a termination signal.
func (s *Server) Run(ctx context.Context) error {
	return server.Serve(ctx, s.srv)
}
