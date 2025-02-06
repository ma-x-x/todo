package middleware

import (
	"time"

	"todo/pkg/monitor"

	"github.com/gin-gonic/gin"
)

func DBMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 记录数据库操作时间
		duration := time.Since(start).Seconds()
		monitor.DatabaseQueryDuration.WithLabelValues("query").Observe(duration)
	}
}
