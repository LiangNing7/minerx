package app

import (
	"fmt"

	"github.com/LiangNing7/goutils/pkg/app"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/LiangNing7/minerx/cmd/minerx-toyblc/app/options"
)

// Define the description of the command.
const commandDesc = `The toy blc is used to start a naive and simple blockchain node.`

// NewApp creates and returns a new App object with default parameters.
func NewApp() *app.App {
	opts := options.NewServerOptions()
	application := app.NewApp(
		"minerx-toyblc",
		"Launch a minerx toy blockchain node",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// run contains the main logic for initializing and running the server.
func run(opts *options.ServerOptions) app.RunFunc {
	return func() error {
		// Load the configuration options
		cfg, err := opts.Config()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		ctx := genericapiserver.SetupSignalContext()

		// Build the server using the configuration
		server, err := cfg.NewServer(ctx)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		// Run the server with signal context for graceful shutdown
		return server.Run(ctx)
	}
}
