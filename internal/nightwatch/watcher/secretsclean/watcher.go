// Package secretsclean is a watcher implement used to delete expired keys from the database.
package secretsclean

import (
	"context"
	"time"

	"github.com/LiangNing7/goutils/pkg/log"
	"github.com/LiangNing7/goutils/pkg/store/where"
	"github.com/LiangNing7/goutils/pkg/watch/registry"

	"github.com/LiangNing7/minerx/internal/nightwatch/watcher"
	"github.com/LiangNing7/minerx/internal/pkg/client/store"
)

var _ registry.Watcher = (*secretsCleanWatcher)(nil)

// watcher implement.
type secretsCleanWatcher struct {
	store store.Interface
}

// Run runs the watcher.
func (w *secretsCleanWatcher) Run() {
	ctx := context.Background()
	_, secrets, err := w.store.UserCenter().Secret().List(ctx, nil)
	if err != nil {
		log.Errorw(err, "Failed to list secrets")
		return
	}

	for _, secret := range secrets {
		if secret.Expires != 0 && secret.Expires < time.Now().AddDate(0, 0, -7).Unix() {
			err := w.store.UserCenter().Secret().Delete(ctx, where.F("user_id", secret.UserID, "name", secret.Name))
			if err != nil {
				log.Warnw("Failed to delete secret from database", "userID", secret.UserID, "name", secret.Name)
				continue
			}
			log.Infow("Successfully deleted secret from database", "userID", secret.UserID, "name", secret.Name)
		}
	}
}

// SetAggregateConfig initializes the watcher for later execution.
func (w *secretsCleanWatcher) SetAggregateConfig(config *watcher.AggregateConfig) {
	w.store = config.AggregateStore
}

func init() {
	registry.Register("secretsclean", &secretsCleanWatcher{})
}
