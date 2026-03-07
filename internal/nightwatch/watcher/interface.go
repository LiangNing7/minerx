package watcher

import (
	"github.com/LiangNing7/goutils/pkg/watch/registry"

	"github.com/LiangNing7/minerx/internal/nightwatch/store"
	aggregatestore "github.com/LiangNing7/minerx/internal/pkg/client/store"
)

// WantsAggregateConfig defines a function which sets AggregateConfig for watcher plugins that need it.
type WantsAggregateConfig interface {
	registry.Watcher
	SetAggregateConfig(config *AggregateConfig)
}

// WantsAggregateStore defines a function which sets aggregate store for watcher plugins that need it.
type WantsAggregateStore interface {
	registry.Watcher
	SetAggregateStore(store aggregatestore.Interface)
}

// WantsStore defines a function which sets store for watcher plugins that need it.
type WantsStore interface {
	registry.Watcher
	SetStore(store store.IStore)
}
