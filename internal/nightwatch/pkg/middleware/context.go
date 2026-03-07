package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/LiangNing7/minerx/internal/pkg/known"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(known.XUserID, c.GetHeader(known.XUserID))
		c.Next()
	}
}
