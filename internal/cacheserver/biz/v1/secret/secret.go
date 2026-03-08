package secret

//go:generate mockgen -destination mock_secret.go -package secret github.com/LiangNing7/minerx/internal/cacheserver/biz/secret SecretBiz

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/LiangNing7/minerx/internal/usercenter/model"
	v1 "github.com/LiangNing7/minerx/pkg/api/cacheserver/v1"
	"github.com/LiangNing7/minerx/pkg/cache"
)

// SecretBiz defines the interface that contains methods for handling secret requests.
type SecretBiz interface {
	Set(ctx context.Context, rq *v1.SetSecretRequest) (*emptypb.Empty, error)
	Del(ctx context.Context, rq *v1.DelSecretRequest) (*emptypb.Empty, error)
	Get(ctx context.Context, rq *v1.GetSecretRequest) (*v1.GetSecretResponse, error)
}

// secretBiz is the implementation of SecretBiz.
type secretBiz struct {
	cache *cache.ChainCache[any]
}

// Ensure that *secretBiz implements the SecretBiz.
var _ SecretBiz = (*secretBiz)(nil)

// New creates and returns a new instance of *secretBiz.
func New(cache *cache.ChainCache[any]) *secretBiz {
	return &secretBiz{cache: cache}
}

// Set stores a secret in the cache.
func (b *secretBiz) Set(ctx context.Context, rq *v1.SetSecretRequest) (*emptypb.Empty, error) {
	secret := &model.SecretM{
		Name:        rq.Name,
		SecretID:    rq.Key,
		Description: rq.Description,
	}
	if rq.Expire != nil {
		secret.Expires = time.Now().Add(rq.Expire.AsDuration()).Unix()
	}

	return &emptypb.Empty{}, b.cache.Set(ctx, rq.Key, secret)
}

// Del deletes a secret from the cache.
func (b *secretBiz) Del(ctx context.Context, rq *v1.DelSecretRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, b.cache.Del(ctx, rq.Key)
}

// Get retrieves a secret from the cache.
func (b *secretBiz) Get(ctx context.Context, rq *v1.GetSecretRequest) (*v1.GetSecretResponse, error) {
	value, err := b.cache.Get(ctx, rq.Key)
	if err != nil {
		return nil, err
	}

	secret := value.(*model.SecretM)

	var rp v1.GetSecretResponse
	_ = copier.Copy(&rp, value)
	rp.CreatedAt = timestamppb.New(secret.CreatedAt)
	rp.UpdatedAt = timestamppb.New(secret.UpdatedAt)
	return &rp, nil
}
