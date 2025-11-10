//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package gateway

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/LiangNing7/goutils/pkg/db"
	"github.com/LiangNing7/goutils/pkg/server"
	genericvalidation "github.com/LiangNing7/goutils/pkg/validation"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/gateway/biz"
	"github.com/LiangNing7/minerx/internal/gateway/handler"
	"github.com/LiangNing7/minerx/internal/gateway/pkg/validation"
	"github.com/LiangNing7/minerx/internal/gateway/store"
	"github.com/LiangNing7/minerx/internal/pkg/client/usercenter"
	"github.com/LiangNing7/minerx/internal/pkg/idempotent"
	"github.com/LiangNing7/minerx/internal/pkg/middleware/validate"
)

func InitializeWebServer(
	<-chan struct{},
	*Config,
	clientset.Interface,
	*db.MySQLOptions,
	*db.RedisOptions,
) (server.Server, error) {
	wire.Build(
		NewWebServer,
		NewMiddlewares,
		ProvideKratosAppConfig,
		wire.NewSet(server.NewEtcdRegistrar, wire.FieldsOf(new(*Config), "EtcdOptions")),
		ProvideKratosLogger,
		handler.ProviderSet,
		store.ProviderSet,
		biz.ProviderSet,
		wire.NewSet(usercenter.ProviderSet, wire.FieldsOf(new(*Config), "UserCenterOptions")),
		db.ProviderSet,
		idempotent.ProviderSet,
		wire.NewSet(
			validation.ProviderSet,
			genericvalidation.NewValidator,
			wire.Bind(new(validate.RequestValidator), new(*genericvalidation.Validator)),
		),
		createInformers,
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
	)

	return nil, nil
}
