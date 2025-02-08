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

// CacheMiddleware 缓存中间件
func CacheMiddleware(rdb *redis.Client, expire time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := generateCacheKey(c)
		ctx := context.Background()

		// 尝试从缓存获取
		if cached, err := rdb.Get(ctx, key).Result(); err == nil {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			rdb.Set(ctx, key, w.body.String(), expire)
			c.Header("X-Cache", "MISS")
		}
	}
}

func generateCacheKey(c *gin.Context) string {
	data := fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
	hash := sha256.Sum256([]byte(data))
	return "cache:" + hex.EncodeToString(hash[:])
}
