package handler

import (
	"github.com/google/wire"

	"github.com/LiangNing7/minerx/internal/usercenter/biz"
	v1 "github.com/LiangNing7/minerx/pkg/api/usercenter/v1"
)

// ProviderSet contains providers for creating instances of the biz struct.
var ProviderSet = wire.NewSet(NewHandler, wire.Bind(new(v1.UserCenterServer), new(*Handler)))

// Handler implements a gRPC service.
type Handler struct {
	v1.UnimplementedUserCenterServer

	biz biz.IBiz
}

// Ensure that Handler implements the v1.UserCenterServer interface.
var _ v1.UserCenterServer = (*Handler)(nil)

// NewHandler creates a new instance of *Handler.
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{biz: biz}
}
