// Package index provides indexes for the api.
package index

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

// AddDefaultIndexes registers the default list of indexes.
func AddDefaultIndexes(ctx context.Context, mgr ctrl.Manager) error {
	if err := ByMinerPod(ctx, mgr); err != nil {
		return err
	}

	return nil
}
