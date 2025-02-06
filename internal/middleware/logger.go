package middleware

import (
	"time"
	"todo/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
// 记录HTTP请求的详细信息,包括请求方法、URI、状态码、客户端IP和响应时间等
// @return gin.HandlerFunc 返回Gin中间件处理函数
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求,调用下一个处理函数
		c.Next()

		// 记录请求结束时间
		endTime := time.Now()

		// 计算请求处理耗时
		latencyTime := endTime.Sub(startTime)

		// 获取请求相关信息
		reqMethod := c.Request.Method   // 请求方法(GET/POST等)
		reqUri := c.Request.RequestURI  // 请求URI
		statusCode := c.Writer.Status() // HTTP状态码
		clientIP := c.ClientIP()        // 客户端IP地址

		// 使用结构化日志记录请求信息
		// 包含请求方法、URI、状态码、客户端IP和处理耗时
		logger.Info().
			Str("method", reqMethod).
			Str("uri", reqUri).
			Int("status", statusCode).
			Str("ip", clientIP).
			Dur("latency", latencyTime).
			Msg("HTTP Request")
	}
}
