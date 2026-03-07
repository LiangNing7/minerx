// Package clean is a watcher implement.
package clean

import (
	"context"

	"github.com/LiangNing7/goutils/pkg/log"
	"github.com/LiangNing7/goutils/pkg/watch/registry"

	"github.com/LiangNing7/minerx/internal/nightwatch/watcher"
	"github.com/LiangNing7/minerx/internal/pkg/client/store"
)

var _ registry.Watcher = (*Watcher)(nil)

// watcher implement.
type Watcher struct {
	store store.Interface
}

// Run runs the watcher.
func (w *Watcher) Run() {
	_, miners, err := w.store.Gateway().Miner().List(context.Background(), nil)
	if err != nil {
		log.Errorw(err, "Failed to list miners")
		return
	}

	for _, m := range miners {
		log.Infow("Retrieve a miner", "miner", m.Name)
	}
}

// SetAggregateConfig initializes the watcher for later execution.
func (w *Watcher) SetAggregateConfig(config *watcher.AggregateConfig) {
	w.store = config.AggregateStore
}

func init() {
	registry.Register("clean", &Watcher{})
}
