package usercenter

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned when a method is not implemented.
var ErrNotImplemented = errors.New("not implemented")

func (c *clientImpl) GetSecret(ctx context.Context, rq *GetSecretRequest) (*GetSecretResponse, error) {
	return nil, ErrNotImplemented
}
