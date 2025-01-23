package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Cache 缓存中间件配置
type CacheConfig struct {
	Expiration time.Duration
	KeyPrefix  string
}

// CacheMiddleware returns a middleware that caches responses in Redis
func CacheMiddleware(rdb *redis.Client, expire time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成缓存key
		key := generateCacheKey(c)

		// 尝试从缓存获取
		ctx := context.Background()
		if cached, err := rdb.Get(ctx, key).Result(); err == nil {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		// 创建自定义的 ResponseWriter
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 如果是成功响应，则缓存
		if c.Writer.Status() == http.StatusOK {
			rdb.Set(ctx, key, w.body.String(), expire)
			c.Header("X-Cache", "MISS")
		}
	}
}

// generateCacheKey generates a unique cache key based on the request
func generateCacheKey(c *gin.Context) string {
	// 组合请求信息
	data := fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)

	// 计算哈希
	hash := sha256.Sum256([]byte(data))
	return "cache:" + hex.EncodeToString(hash[:])
}
