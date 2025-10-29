package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/LiangNing7/minerx/internal/pkg/contextx"
	known "github.com/LiangNing7/minerx/internal/pkg/known/toyblc"
)

// TraceID is a Gin middleware to inject the `Trace-ID` key-value pair into the HTTP request context and response.
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the `Trace-ID` key-value pair is in the HTTP request header, if it is, reuse it, otherwise create a new one.
		traceID := c.Request.Header.Get(known.TraceIDKey)

		if traceID == "" {
			traceID = uuid.New().String()
			// Set the `Trace-ID` key-value pair into the HTTP request header.
			c.Request.Header.Set(known.TraceIDKey, traceID)
		}

		// Set the `Trace-ID` key-value pair into the HTTP response header.
		c.Writer.Header().Set(known.TraceIDKey, traceID)

		// Set the `trace.id` key-value pair into the Gin context.
		// Use `trace.id` instead of `known.TraceIDKey` to keep consistent with other components.
		c.Set("trace.id", traceID)

		ctx := contextx.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
