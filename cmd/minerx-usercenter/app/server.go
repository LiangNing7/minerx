package app

import (
	"context"
	"fmt"

	"github.com/LiangNing7/goutils/pkg/app"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/LiangNing7/minerx/cmd/minerx-usercenter/app/options"
	"github.com/LiangNing7/minerx/internal/pkg/contextx"
	"github.com/LiangNing7/minerx/internal/pkg/known"
	"github.com/LiangNing7/minerx/internal/usercenter"
)

// Define the description of the command.
const commandDesc = `The usercenter server is used to manage users, keys, fees, etc.`

// NewApp creates and returns a new App object with default parameters.
func NewApp() *app.App {
	opts := options.NewServerOptions()
	application := app.NewApp(
		usercenter.Name,
		"Launch a minerx usercenter server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithLoggerContextExtractor(map[string]func(context.Context) string{
			known.XTraceID: contextx.TraceID,
			known.XUserID:  contextx.UserID,
		}),
	)

	return application
}

func run(opts *options.ServerOptions) app.RunFunc {
	return func() error {
		// Load the configuration options.
		cfg, err := opts.Config()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		ctx := genericapiserver.SetupSignalContext()

		// Build the server using the configuration.
		server, err := cfg.NewServer(ctx)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		// Run the server with signal context for graceful shutdown.
		return server.Run(ctx)
	}
}
