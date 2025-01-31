package handlers

import (
	"net/http"
	"todo-demo/api/v1/dto/auth"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Register 用户注册处理器
// @Summary 用户注册
// @Description 创建新用户账号，验证注册信息并在成功时返回确认消息
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册信息，包含用户名和密码"
// @Success 200 {object} auth.RegisterResponse "注册成功返回信息"
// @Failure 400 {object} errors.Error "请求参数错误或注册过程中的业务错误"
// @Router /auth/register [post]
func Register(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := authService.Register(c.Request.Context(), &req); err != nil {
			if err == errors.ErrUserExists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}

		c.JSON(http.StatusOK, auth.RegisterResponse{
			Message: "注册成功",
		})
	}
}

// Login 用户登录处理器
// @Summary 用户登录
// @Description 验证用户凭证并生成JWT令牌，用于后续请求的身份验证
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录信息，包含用户名和密码"
// @Success 200 {object} auth.LoginResponse "登录成功返回的JWT令牌"
// @Failure 401 {object} errors.Error "用户名或密码错误等认证失败的情况"
// @Router /auth/login [post]
func Login(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := authService.Login(c.Request.Context(), &req)
		if err != nil {
			if err == errors.ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
			return
		}

		c.JSON(http.StatusOK, auth.LoginResponse{
			Token: token,
		})
	}
}
