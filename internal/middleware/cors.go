package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域资源共享(CORS)中间件
// 设置跨域相关的HTTP响应头,允许跨域请求访问API
// @return gin.HandlerFunc 返回Gin中间件处理函数
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源的跨域请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许客户端发送身份凭证(cookies等)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 允许的HTTP请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 对于OPTIONS预检请求,直接返回204状态码
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
