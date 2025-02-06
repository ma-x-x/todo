package middleware

import (
	"time"

	"todo/pkg/monitor"

	"github.com/gin-gonic/gin"
)

func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 请求计数
		monitor.RequestCounter.WithLabelValues(method, path, "200").Inc()

		// 执行请求
		c.Next()

		// 记录请求时间
		duration := time.Since(start).Seconds()
		monitor.RequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
