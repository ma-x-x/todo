package middleware

import (
	"bytes"
	"io"
	"time"
	"todo/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
// 记录HTTP请求的详细信息,包括请求方法、URI、状态码、客户端IP和响应时间等
// @return gin.HandlerFunc 返回Gin中间件处理函数
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContext()
		SetContext(c, ctx)

		start := time.Now()

		// 记录请求体大小
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		latency := time.Since(start)

		// 使用结构化日志记录请求信息
		fields := map[string]interface{}{
			"client_ip":    c.ClientIP(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"status":       c.Writer.Status(),
			"latency":      latency,
			"user_id":      ctx.UserID,
			"request_size": len(requestBody),
			"headers":      c.Request.Header,
		}

		logger.LogInfo(
			ctx.TraceID,
			"HTTP Request",
			fields,
		)
	}
}
