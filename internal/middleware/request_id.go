package middleware

import (
	"github.com/gin-gonic/gin"
)

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := GetContext(c)
		
		// 尝试从请求头获取请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = ctx.RequestID
		}
		
		// 设置请求ID到响应头
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
} 