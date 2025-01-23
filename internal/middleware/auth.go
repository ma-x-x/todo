package middleware

import (
	"strings"
	"todo-demo/pkg/config"
	"todo-demo/pkg/errors"
	"todo-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, errors.NewError(401, errors.ErrUnauthorized.Error()))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, errors.NewError(401, errors.ErrInvalidToken.Error()))
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg)
		if err != nil {
			c.AbortWithStatusJSON(401, errors.NewError(401, err.Error()))
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
