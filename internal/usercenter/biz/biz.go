package biz

//go:generate mockgen -destination mock_biz.go -package biz minerx/internal/usercenter/biz IBiz

import (
	"github.com/LiangNing7/goutils/pkg/authn"
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/usercenter/pkg/auth"
	"github.com/LiangNing7/minerx/internal/usercenter/store"
)

// ProviderSet is a Wire provider set used to declare dependency injection rules.
// Includes the NewBiz constructor to create a biz instance.
// wire.Bind binds the IBiz interface to the concrete implementation *biz,
// so places that depend on IBiz will automatically inject a *biz instance.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// IBiz defines the methods that must be implemented by the business layer.
type IBiz interface{}

// biz is a concrete implementation of IBiz.
type biz struct {
	store store.IStore
	authn authn.Authenticator
	auth  auth.AuthProvider
}

// Ensure that biz implements the IBiz.
var _ IBiz = (*biz)(nil)

// NewBiz creates an instance of IBiz.
func NewBiz(store store.IStore, authn authn.Authenticator, auth auth.AuthProvider) *biz {
	return &biz{store: store, authn: authn, auth: auth}
}
