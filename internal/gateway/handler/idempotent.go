package handler

import (
	"context"

	v1 "github.com/LiangNing7/minerx/pkg/api/gateway/v1"
)

func (s *Handler) GetIdempotentToken(ctx context.Context, rq *v1.IdempotentRequest) (*v1.IdempotentResponse, error) {
	return &v1.IdempotentResponse{Token: s.idt.Token(ctx)}, nil
}
