package middleware

import (
	"strings"
	"todo-demo/pkg/config"
	"todo-demo/pkg/errors"
	"todo-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
// 用于验证请求头中的JWT令牌,确保API的安全访问
// @param cfg *config.JWTConfig JWT配置信息,包含密钥等配置
// @return gin.HandlerFunc 返回Gin中间件处理函数
func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 如果请求头为空,返回401未授权错误
			c.AbortWithStatusJSON(401, errors.NewError(401, errors.ErrUnauthorized.Error()))
			return
		}

		// 解析Authorization头,格式必须为: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// 如果格式不正确,返回401无效令牌错误
			c.AbortWithStatusJSON(401, errors.NewError(401, errors.ErrInvalidToken.Error()))
			return
		}

		// 解析JWT令牌
		claims, err := utils.ParseToken(parts[1], cfg)
		if err != nil {
			// 如果令牌解析失败,返回401错误
			c.AbortWithStatusJSON(401, errors.NewError(401, err.Error()))
			return
		}

		// 将用户ID保存到上下文中,供后续处理函数使用
		c.Set("userID", claims.UserID)
		// 调用下一个处理函数
		c.Next()
	}
}
