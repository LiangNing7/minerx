package registry

import (
	"fmt"

	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// REST implements a RESTStorage for API services against etcd.
type REST struct {
	*genericregistry.Store
}

// RESTInPeace is just a simple function that panics on error.
// Otherwise returns the given storage object. It is meant to be
// a wrapper for minerx server registries.
func RESTInPeace(storage rest.StandardStorage, err error) rest.StandardStorage {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource due to %w, will die", err)
		panic(err)
	}
	return storage
}
