package auth

// LoginRequest 登录请求
type LoginRequest struct {
	// Username 用户名
	Username string `json:"username" binding:"required"`
	// Password 密码
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
