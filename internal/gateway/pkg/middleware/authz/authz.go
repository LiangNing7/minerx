package auth

import (
	"context"

	"github.com/LiangNing7/goutils/pkg/i18n"
	"github.com/LiangNing7/goutils/pkg/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/LiangNing7/minerx/internal/gateway/pkg/locales"
	"github.com/LiangNing7/minerx/internal/pkg/contextx"
	"github.com/LiangNing7/minerx/internal/pkg/middleware/authz"
	jwtutil "github.com/LiangNing7/minerx/internal/pkg/util/jwt"
	"github.com/LiangNing7/minerx/pkg/api/errno"
)

// Authz is a authentication and authorization middleware.
func Authz(authz authz.Authorizer) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq any) (reply any, err error) {
			accessToken := jwtutil.TokenFromServerContext(ctx)
			if tr, ok := transport.FromServerContext(ctx); ok {
				userID, allowed, err := authz.Authorize(ctx, accessToken, "*", tr.Operation())
				if err != nil {
					log.Errorw(err, "Authorization failure occurs", "operation", tr.Operation())
					return nil, err
				}
				if !allowed {
					return nil, errno.ErrorForbidden("%s", i18n.FromContext(ctx).T(locales.NoPermission))
				}
				ctx = contextx.WithUserID(ctx, userID)
				ctx = contextx.WithAccessToken(ctx, accessToken)
				ctx = contextx.WithUserID(ctx, userID)
			}

			return handler(ctx, rq)
		}
	}
}
