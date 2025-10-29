//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package usercenter

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/LiangNing7/goutils/pkg/db"
	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	"github.com/LiangNing7/goutils/pkg/server"
	genericvalidation "github.com/LiangNing7/goutils/pkg/validation"
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/pkg/middleware/validate"
	"github.com/LiangNing7/minerx/internal/usercenter/biz"
	"github.com/LiangNing7/minerx/internal/usercenter/handler"
	"github.com/LiangNing7/minerx/internal/usercenter/pkg/auth"
	"github.com/LiangNing7/minerx/internal/usercenter/pkg/validation"
	"github.com/LiangNing7/minerx/internal/usercenter/store"
)

func InitializeWebServer(
	<-chan struct{},
	*Config,
	*db.MySQLOptions,
	*genericoptions.JWTOptions,
	*genericoptions.RedisOptions,
	*genericoptions.KafkaOptions,
) (server.Server, error) {
	wire.Build(
		wire.NewSet(server.NewEtcdRegistrar, wire.FieldsOf(new(*Config), "EtcdOptions")), // dep by AppConfig
		ProvideKratosAppConfig, // server.KratosAppConfig, dep by NewWebServer
		ProvideKratosLogger,    // dep by NewMiddlewares
		// func NewMiddlewares(logger krtlog.Logger,authn authn.Authenticator, val validate.RequestValidator) []middleware.Middleware {
		NewAuthenticator,
		NewWebServer,
		NewMiddlewares,
		store.SetterProviderSet,
		auth.ProviderSet,
		handler.ProviderSet,
		store.ProviderSet,
		biz.ProviderSet,
		db.ProviderSet,
		wire.NewSet(
			validation.ProviderSet,
			genericvalidation.NewValidator,
			wire.Bind(new(validate.RequestValidator), new(*genericvalidation.Validator)),
		),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
	)
	return nil, nil
}
