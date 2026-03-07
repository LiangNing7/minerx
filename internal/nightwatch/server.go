package nightwatch

import (
	"context"

	"github.com/LiangNing7/goutils/pkg/core"
	"github.com/LiangNing7/goutils/pkg/db"
	"github.com/LiangNing7/goutils/pkg/log"
	genericoptions "github.com/LiangNing7/goutils/pkg/options"
	"github.com/LiangNing7/goutils/pkg/store/where"
	"github.com/LiangNing7/goutils/pkg/watch"
	"github.com/LiangNing7/goutils/pkg/watch/logger/onex"
	"gorm.io/gorm"

	"github.com/LiangNing7/minerx/internal/nightwatch/watcher"
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/all"
	fakeminio "github.com/LiangNing7/minerx/internal/pkg/client/minio/fake"
	"github.com/LiangNing7/minerx/internal/pkg/contextx"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
)

// Config contains application-related configurations.
type Config struct {
	HealthOptions     *genericoptions.HealthOptions
	MySQLOptions      *genericoptions.MySQLOptions
	RedisOptions      *genericoptions.RedisOptions
	WatchOptions      *watch.Options
	HTTPOptions       *genericoptions.HTTPOptions
	TLSOptions        *genericoptions.TLSOptions
	DisableRESTServer bool
	// The maximum concurrency event of user watcher.
	UserWatcherMaxWorkers int64
	// The list of watchers that should be disabled.
	Client clientset.Interface
	// Created from MySQLOptions.
	DB *gorm.DB
}

// Server represents the web server.
type Server struct {
	watch *watch.Watch
}

// NewServer initializes and returns a new Server instance.
func (cfg *Config) NewServer(ctx context.Context) (*Server, error) {
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	var mysqlOptions db.MySQLOptions
	_ = core.Copy(&mysqlOptions, cfg.MySQLOptions)
	dbIns, err := db.NewMySQL(&mysqlOptions)
	if err != nil {
		return nil, err
	}
	cfg.DB = dbIns

	watcherConfig, err := cfg.CreateWatcherConfig()
	if err != nil {
		return nil, err
	}

	initialize := watcher.NewInitializer(watcherConfig)
	opts := []watch.Option{
		watch.WithInitialize(initialize),
		watch.WithLogger(onex.NewLogger()),
	}

	watchIns, err := watch.NewWatch(cfg.WatchOptions, cfg.DB, opts...)
	if err != nil {
		return nil, err
	}

	if !cfg.DisableRESTServer {
		go NewRESTServer(cfg.HTTPOptions, cfg.TLSOptions, cfg.DB).RunOrDie()
	} else {
		go cfg.HealthOptions.ServeHealthCheck()
	}

	return &Server{watch: watchIns}, nil
}

// CreateWatcherConfig used to create configuration used by all watcher.
func (cfg *Config) CreateWatcherConfig() (*watcher.AggregateConfig, error) {
	storeClient, err := InitializeStore(cfg.DB)
	if err != nil {
		log.Errorw(err, "Failed to create MySQL client")
		return nil, err
	}

	aggregateStoreClient, err := InitializeAggregateStore(cfg.DB)
	if err != nil {
		log.Errorw(err, "Failed to create MySQL client")
		return nil, err
	}

	minioClient, err := fakeminio.NewFakeMinioClient("test-bucket-name")
	if err != nil {
		log.Errorw(err, "Failed to NewMinioClient")
		return nil, err
	}
	return &watcher.AggregateConfig{
		Minio:                 minioClient,
		Store:                 storeClient,
		AggregateStore:        aggregateStoreClient,
		Client:                cfg.Client,
		UserWatcherMaxWorkers: cfg.UserWatcherMaxWorkers,
	}, nil
}

// Run starts the server and listens for termination signals.
// It gracefully shuts down the server upon receiving a termination signal.
func (s *Server) Run(ctx context.Context) error {
	s.watch.Start(ctx.Done())

	// Wait for stop signal
	<-ctx.Done()
	s.watch.Stop()
	return nil
}
