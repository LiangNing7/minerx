package handler

import (
	"strings"

	"github.com/LiangNing7/goutils/pkg/core"
	"github.com/gin-gonic/gin"

	"github.com/LiangNing7/minerx/internal/toyblc/pkg/ws"
	v1 "github.com/LiangNing7/minerx/pkg/api/toyblc/v1"
)

// CreatePeer handles the creation of a new peer.
func (h *Handler) CreatePeer(c *gin.Context) {
	var rq v1.CreatePeerRequest
	if err := core.ShouldBindJSON(c, &rq); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	ws.ConnectToPeers(c, h.bs, h.ss, []string{rq.Peer})

	core.WriteResponse(c, nil, nil)
}

// ListPeer retrieves a list of peers based on query parameters.
func (h *Handler) ListPeer(c *gin.Context) {
	var slice []string

	for _, socket := range h.ss.List() {
		if socket.IsClientConn() {
			slice = append(slice, strings.Replace(socket.LocalAddr().String(), "ws://", "", 1))
		} else {
			slice = append(slice, socket.Request().RemoteAddr)
		}
	}

	core.WriteResponse(c, slice, nil)
}
