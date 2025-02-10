package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		log.Printf("收到请求: Method=%s, Origin=%s, Path=%s",
			c.Request.Method, origin, c.Request.URL.Path)

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			log.Printf("处理预检请求: Origin=%s", origin)
			c.Writer.WriteHeader(http.StatusOK)
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
