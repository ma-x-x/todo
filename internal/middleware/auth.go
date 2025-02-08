package middleware

import (
	"strings"
	"todo/pkg/config"
	"todo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
// 用于验证请求头中的JWT令牌,确保API的安全访问
//
// Parameters:
//   - jwtCfg: JWT配置信息,包含密钥等配置
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
// @Failure 401 {object} gin.H "未授权访问"
// @Router /auth/middleware [get]
func Auth(jwtCfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "未授权访问",
			})
			return
		}

		claims, err := utils.ParseToken(token, jwtCfg)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "无效的令牌",
			})
			return
		}

		// 将用户ID存储到上下文中
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// extractToken 从请求头中提取 token
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}
