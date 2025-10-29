package authz

import "context"

type Authorizer interface {
	Authorize(ctx context.Context, token string, obj, act string) (userID string, allowed bool, err error)
}
