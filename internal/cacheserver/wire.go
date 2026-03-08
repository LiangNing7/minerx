//go:build wireinject
// +build wireinject

package cacheserver

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/LiangNing7/goutils/pkg/db"
	"github.com/LiangNing7/goutils/pkg/server"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/cacheserver/biz"
	"github.com/LiangNing7/minerx/internal/cacheserver/handler"
	"github.com/LiangNing7/minerx/internal/cacheserver/store"
	"github.com/LiangNing7/minerx/pkg/cache"
)

func InitializeWebServer(*Config, *db.MySQLOptions, cache.Cache[*any.Any], bool) (server.Server, error) {
	wire.Build(
		NewWebServer,
		db.ProviderSet,
		store.ProviderSet,
		biz.ProviderSet,
		handler.ProviderSet,
		wire.Struct(new(ServerConfig), "*"),
	)
	return nil, nil
}
