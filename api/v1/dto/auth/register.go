package auth

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
