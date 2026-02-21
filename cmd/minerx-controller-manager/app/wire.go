//go:build wireinject
// +build wireinject

package app

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/google/wire"

	"github.com/LiangNing7/goutils/pkg/db"

	"github.com/LiangNing7/minerx/internal/gateway/store"
)

func wireStoreClient(*db.MySQLOptions) (store.IStore, error) {
	wire.Build(
		db.ProviderSet,
		store.ProviderSet,
	)

	return nil, nil
}
