//go:build wireinject
// +build wireinject

package nightwatch

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	gwstore "github.com/LiangNing7/minerx/internal/gateway/store"
	"github.com/LiangNing7/minerx/internal/nightwatch/biz"
	"github.com/LiangNing7/minerx/internal/nightwatch/handler"
	"github.com/LiangNing7/minerx/internal/nightwatch/pkg/validation"
	nwstore "github.com/LiangNing7/minerx/internal/nightwatch/store"
	"github.com/LiangNing7/minerx/internal/pkg/client/store"
	ucstore "github.com/LiangNing7/minerx/internal/usercenter/store"
)

func InitializeAggregateStore(*gorm.DB) (store.Interface, error) {
	wire.Build(
		store.ProviderSet,
		gwstore.ProviderSet,
		ucstore.ProviderSet,
	)

	return nil, nil
}

func InitializeHandler(*gorm.DB) *handler.Handler {
	wire.Build(
		validation.New,
		biz.ProviderSet,
		nwstore.ProviderSet,
		handler.NewHandler,
	)

	return nil
}

func InitializeStore(*gorm.DB) (nwstore.IStore, error) {
	wire.Build(nwstore.ProviderSet)

	return nil, nil
}
