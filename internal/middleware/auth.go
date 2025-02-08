package middleware

import (
	"net/http"
	"strings"
	"todo/pkg/config"
	"todo/pkg/utils"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
// 用于验证请求头中的JWT令牌,确保API的安全访问
//
// Parameters:
//   - cfg: JWT配置信息,包含密钥等配置
//
// Returns:
//   - gin.HandlerFunc: 返回Gin中间件处理函数
//
// @Summary JWT认证中间件
// @Description 用于验证请求头中的JWT令牌,确保API的安全访问
// @Tags middleware
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Success 200 {object} interface{} "验证成功"
// @Failure 401 {object} errors.Error "未授权访问"
// @Router /auth/middleware [get]
func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				http.StatusUnauthorized, "未经授权的访问"))
			return
		}

		claims, err := utils.ParseToken(token, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				http.StatusUnauthorized, "无效的访问令牌"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
