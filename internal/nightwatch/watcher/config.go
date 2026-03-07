// Package watcher provides functions used by all watchers.
package watcher

import (
	"github.com/LiangNing7/minerx/internal/nightwatch/store"
	"github.com/LiangNing7/minerx/internal/pkg/client/minio"
	aggregatestore "github.com/LiangNing7/minerx/internal/pkg/client/store"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
)

// AggregateConfig aggregates the configurations of all watchers and serves as a configuration aggregator.
type AggregateConfig struct {
	Minio minio.IMinio
	Store store.IStore
	// The purpose of nightwatch is to handle asynchronous tasks on the minerx platform
	// in a unified manner, so a store aggregation type is needed here.
	AggregateStore aggregatestore.Interface

	// Client is the client for minerx-apiserver.
	Client clientset.Interface

	// Then maximum concurrency event of user watcher.
	UserWatcherMaxWorkers int64
}
