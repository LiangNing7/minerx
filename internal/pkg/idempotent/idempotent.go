package idempotent

import (
	"context"

	"github.com/LiangNing7/goutils/pkg/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"github.com/LiangNing7/minerx/pkg/idempotent"
)

var ProviderSet = wire.NewSet(NewIdempotent)

type Idempotent struct {
	idempotent *idempotent.Idempotent
}

func (idt Idempotent) Token(ctx context.Context) string {
	return idt.idempotent.Token(ctx)
}

func (idt Idempotent) Check(ctx context.Context, token string) bool {
	return idt.idempotent.Check(ctx, token)
}

// NewIdempotent is initialize idempotent from config.
func NewIdempotent(redis redis.UniversalClient) (idt *Idempotent, err error) {
	ins := idempotent.New(idempotent.WithRedis(redis))
	idt = &Idempotent{
		idempotent: ins,
	}

	log.Infow("Initialize idempotent success")
	return idt, nil
}
