package middleware

import (
	"github.com/gin-gonic/gin"
)

// TraceMiddleware 追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := GetContext(c)
		
		// 设置追踪ID到响应头
		c.Header("X-Trace-ID", ctx.TraceID)
		c.Header("X-Request-ID", ctx.RequestID)
		
		c.Next()
	}
} 