package validation

import (
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/gateway/store"
)

// Validator is a struct that implements custom validation logic.
type Validator struct {
	// Some complex validation logic may require direct database queries.
	// This is just an example. If validation requires other dependencies
	// like clients, services, resources, etc., they can all be injected here.
	store store.IStore
}

// ProviderSet is a Wire provider set that declares dependency injection rules.
var ProviderSet = wire.NewSet(New, wire.Bind(new(any), new(*Validator)))

// New creates a new instance of Validator.
func New(store store.IStore) *Validator {
	return &Validator{store: store}
}
