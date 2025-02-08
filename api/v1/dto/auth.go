// Package dto 提供数据传输对象定义
package dto

// LoginRequest 登录请求参数
type LoginRequest struct {
	// Username 用户名
	// Required: true
	// Min length: 3
	// Max length: 32
	Username string `json:"username" binding:"required,min=3,max=32"`

	// Password 密码
	// Required: true
	// Min length: 6
	// Max length: 32
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	// JWT令牌
	Token string `json:"token"`

	// 用户信息
	User *UserInfo `json:"user"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	// Username 用户名
	Username string `json:"username" binding:"required,min=3,max=32"`
	// Password 密码
	Password string `json:"password" binding:"required,min=6,max=32"`
	// Email 邮箱
	Email string `json:"email" binding:"required,email"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	// Message 响应消息
	Message string `json:"message"`
}

// UserInfo 用户信息
type UserInfo struct {
	// 用户ID
	ID uint `json:"id"`

	// 用户名
	Username string `json:"username"`

	// 邮箱
	Email string `json:"email"`

	// 创建时间
	CreatedAt string `json:"createdAt"`

	// 更新时间
	UpdatedAt string `json:"updatedAt"`
}
