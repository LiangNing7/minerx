package cacheserver

import (
	"github.com/LiangNing7/goutils/pkg/server"
	"google.golang.org/grpc"

	v1 "github.com/LiangNing7/minerx/pkg/api/cacheserver/v1"
)

func (c *ServerConfig) NewGRPCServer() (server.Server, error) {
	// Initialize the gRPC server.
	return server.NewGRPCServer(
		c.cfg.GRPCOptions,
		c.cfg.TLSOptions,
		[]grpc.ServerOption{},
		func(s grpc.ServiceRegistrar) { v1.RegisterCacheServerServer(s, c.handler) })
}
