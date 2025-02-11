package handlers

import (
	"net/http"
	"todo/api/v1/dto/auth"
	"todo/internal/service"
	"todo/pkg/errors"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=auth.RegisterResponse} "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层处理注册逻辑
	err := h.authService.Register(c.Request.Context(), &req)
	if err == errors.ErrUserExists {
		response.Error(c, http.StatusBadRequest, "用户名已存在")
		return
	} else if err != nil {
		response.Error(c, http.StatusInternalServerError, "注册失败")
		return
	}

	response.Success(c, auth.RegisterResponse{
		Message: "注册成功",
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取JWT令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=gin.H{token=string,user=auth.UserInfo}} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 优化验证错误提示
		if validationErr := errors.ParseValidationError(err); validationErr != "" {
			response.BadRequest(c, validationErr)
			return
		}
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务层处理登录逻辑
	token, userInfo, err := h.authService.Login(c.Request.Context(), &req)
	if err == errors.ErrInvalidCredentials {
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	} else if err != nil {
		response.Error(c, http.StatusInternalServerError, "登录失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  userInfo,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 注销用户登录状态
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Success 200 {object} response.Response "登出成功"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取用户ID
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	// 调用服务层处理登出逻辑
	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "登出失败")
		return
	}

	response.Success(c, nil)
}
