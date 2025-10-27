package usercenter

import (
	"context"
	"os"

	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	"github.com/LiangNing7/goutils/pkg/server"
	"github.com/LiangNing7/goutils/pkg/version"
)

var (
	// Name is name of the compiled software.
	Name = "minerx-usercenter"

	ID, _ = os.Hostname()

	Version = version.Get().String()
)

type Config struct {
	GRPCOptions   *genericoptions.GRPCOptions
	HTTPOptions   *genericoptions.HTTPOptions
	TLSOptions    *genericoptions.TLSOptions
	JWTOptions    *genericoptions.JWTOptions
	MySQLOptions  *genericoptions.MySQLOptions
	RedisOptions  *genericoptions.RedisOptions
	EtcdOptions   *genericoptions.EtcdOptions
	KafkaOptions  *genericoptions.KafkaOptions
	JaegerOptions *genericoptions.JaegerOptions
	ConsulOptions *genericoptions.ConsulOptions
}

// Server represents the web server.
type Server struct {
	srv server.Server
}

// ServerConfig contains the core dependencies and configurations of the server.
type ServerConfig struct {
	cfg *Config
}

// NewServer initializes and returns a new Server instance.
func (cfg *Config) NewServer(ctx context.Context) (*Server, error) {
	// TODO:
	return &Server{}, nil
}

// Run starts the server and listens for termination signals.
// It gracefully shuts down the server upon receiving a termination signal.
func (s *Server) Run(ctx context.Context) error {
	return server.Serve(ctx, s.srv)
}
