package biz

//go:generate mockgen -destination mock_biz.go -package biz github.com/LiangNing7/minerx/internal/gateway/biz IBiz

import (
	"github.com/LiangNing7/minerx/internal/gateway/store"
	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
	"github.com/LiangNing7/minerx/pkg/generated/informers"
	"github.com/google/wire"
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
	store     store.IStore
	clientset clientset.Interface
	informer  informers.SharedInformerFactory
}

// Ensure that biz implements the IBiz.
var _ IBiz = (*biz)(nil)

// NewBiz creates an instance of IBiz.
func NewBiz(store store.IStore, clientset clientset.Interface, informer informers.SharedInformerFactory) *biz {
	return &biz{store: store, clientset: clientset, informer: informer}
}
