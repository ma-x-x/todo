package middleware

import (
	"runtime/debug"
	"todo/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := GetContext(c)

				fields := map[string]interface{}{
					"error":      err,
					"stacktrace": string(debug.Stack()),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}

				logger.LogError(
					ctx.TraceID,
					nil,
					"Panic recovered",
					fields,
				)

				c.AbortWithStatusJSON(500, gin.H{
					"error":    "Internal Server Error",
					"trace_id": ctx.TraceID,
				})
			}
		}()
		c.Next()
	}
}
