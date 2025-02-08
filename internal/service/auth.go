package service

import (
	"context"
	"todo/api/v1/dto/auth"
)

// AuthService 定义认证相关的业务接口
type AuthService interface {
	// Register 用户注册
	// ctx: 上下文信息
	// req: 注册请求，包含用户名、密码等信息
	// 返回错误信息，如果注册成功则返回nil
	Register(ctx context.Context, req *auth.RegisterRequest) error

	// Login 用户登录
	// ctx: 上下文信息
	// req: 登录请求，包含用户名和密码
	// 返回JWT令牌、用户信息和可能的错误
	Login(ctx context.Context, req *auth.LoginRequest) (string, *auth.UserInfo, error)

	// Logout 用户登出
	// ctx: 上下文信息
	// userID: 用户ID
	// 返回可能的错误
	Logout(ctx context.Context, userID uint) error
}
