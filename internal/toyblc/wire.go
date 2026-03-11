//go:build wireinject
// +build wireinject

package toyblc

import (
	"github.com/LiangNing7/goutils/pkg/server"
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/toyblc/pkg/blc"
	"github.com/LiangNing7/minerx/internal/toyblc/pkg/ws"
)

// InitializeWebServer sets up and initializes the web server with all necessary dependencies.
func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		NewAggregatorServer,
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.FieldsOf(new(*Config), "Miner", "MinMineInterval", "Accounts", "Peers"),
		ws.ProviderSet,
		blc.NewBlockSet,
		wire.FieldsOf(new(*Config), "Address"),
	)
	return nil, nil
}
