package jwt

import (
	"context"
	"fmt"
	"strings"

	"github.com/LiangNing7/goutils/pkg/authn"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/LiangNing7/minerx/internal/pkg/contextx"
)

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	// bearerWord the bearer key word for authorization.
	bearerWord string = "Bearer"

	// bearerFormat authorization token format.
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the rquest tokenHeader.
	authorizationKey string = "Authorization"
)

var (
	ErrMissingJwtToken = errors.Unauthorized(reason, "JWT token is missing")
	ErrWrongContext    = errors.Unauthorized(reason, "Wrong context for middleware")
)

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(a authn.Authenticator) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq any) (any, error) {
			// get the transport from the context.
			if tr, ok := transport.FromServerContext(ctx); ok {
				// split the authorization header into two parts.
				auths := strings.SplitN(tr.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}

				accessToken := auths[1]
				claims, err := a.ParseClaims(ctx, accessToken)
				if err != nil {
					return nil, err
				}

				// set the claims to the context.
				ctx = contextx.WithClaims(ctx, claims)
				// set the user id to the context.
				ctx = contextx.WithUserID(ctx, claims.Subject)
				// set the access token to the context.
				ctx = contextx.WithAccessToken(ctx, accessToken)
				// call the next handler.
				return handler(ctx, rq)
			}
			return nil, ErrWrongContext
		}
	}
}

// WithSignToken is a client jwt middleware used to sign a jwt token.
func WithSignToken(a authn.Authenticator, userID string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq any) (any, error) {
			accessToken, err := a.Sign(ctx, userID)
			if err != nil {
				return nil, err
			}

			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(authorizationKey, fmt.Sprintf(bearerFormat, accessToken))
				return handler(ctx, rq)
			}
			return nil, ErrWrongContext
		}
	}
}

// WithToken is a client jwt middleware.
func WithToken(token string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq any) (any, error) {
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(authorizationKey, fmt.Sprintf(bearerFormat, token))
				return handler(ctx, rq)
			}
			return nil, ErrWrongContext
		}
	}
}
