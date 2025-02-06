package handlers

import (
	"net/http"
	"todo-demo/api/v1/dto/auth"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"
	"todo-demo/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 用户注册处理器
// @Summary 用户注册
// @Description 创建新用户账号，验证注册信息并在成功时返回确认消息
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册信息，包含用户名和密码"
// @Success 200 {object} response.Response{data=auth.RegisterResponse} "注册成功返回信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/register [post]
func Register(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
			return
		}

		if err := authService.Register(c.Request.Context(), &req); err != nil {
			if err == errors.ErrUserExists {
				c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "用户名已存在"))
				return
			}
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "注册失败"))
			return
		}

		c.JSON(http.StatusOK, response.Success(auth.RegisterResponse{
			Message: "注册成功",
		}))
	}
}

// Login 用户登录处理器
// @Summary 用户登录
// @Description 验证用户凭证并生成JWT令牌，返回令牌和用户信息
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录信息，包含用户名和密码"
// @Success 200 {object} response.Response{data=auth.LoginResponse} "登录成功返回的JWT令牌和用户信息"
// @Failure 401 {object} response.Response "用户名或密码错误等认证失败的情况"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/login [post]
func Login(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
			return
		}

		token, userInfo, err := authService.Login(c.Request.Context(), &req)
		if err != nil {
			if err == errors.ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "用户名或密码错误"))
				return
			}
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "登录失败"))
			return
		}

		c.JSON(http.StatusOK, response.Success(auth.LoginResponse{
			Token: token,
			User:  userInfo,
		}))
	}
}
