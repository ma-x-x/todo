package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter returns a middleware that limits request rate using Redis
func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		key := "ratelimit:" + ip

		// 使用 Redis 实现简单的计数器限流
		count, err := rdb.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "限流服务错误",
			})
			return
		}

		// 设置过期时间
		if count == 1 {
			rdb.Expire(c, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

// RedisRateLimiter returns a middleware that uses Redis for rate limiting
func RedisRateLimiter(rdb *redis.Client, key string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		now := time.Now().UnixNano()

		// 清理旧的请求记录
		rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window.Nanoseconds()))

		// 添加当前请求
		rdb.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})

		// 获取当前时间窗口内的请求数
		count, err := rdb.ZCard(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "限流服务异常"})
			c.Abort()
			return
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 设置过期时间
		rdb.Expire(ctx, key, window)

		c.Next()
	}
}
