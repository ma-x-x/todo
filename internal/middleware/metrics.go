package middleware

import (
	"time"
	"todo/pkg/monitor"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 指标收集中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 请求计数
		monitor.RequestCounter.WithLabelValues(method, path, "200").Inc()

		c.Next()

		// 记录请求时间和数据库操作时间
		duration := time.Since(start).Seconds()
		monitor.RequestDuration.WithLabelValues(method, path).Observe(duration)
		monitor.DatabaseQueryDuration.WithLabelValues("query").Observe(duration)
	}
}
