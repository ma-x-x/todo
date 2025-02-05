package auth

// LoginRequest 登录请求
type LoginRequest struct {
	// Username 用户名
	Username string `json:"username" binding:"required"`
	// Password 密码
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`              // JWT令牌
	User  *UserInfo   `json:"user"`              // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	ID        uint      `json:"id"`                    // 用户ID
	Username  string    `json:"username"`              // 用户名
	Email     string    `json:"email"`                 // 邮箱
	CreatedAt string    `json:"createdAt"`            // 创建时间
	UpdatedAt string    `json:"updatedAt"`            // 更新时间
}
