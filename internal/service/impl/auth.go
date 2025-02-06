package impl

import (
	"context"
	"todo/api/v1/dto/auth"
	"todo/internal/models"
	"todo/internal/repository"
	"todo/pkg/config"
	"todo/pkg/errors"
	"todo/pkg/utils"
)

// authService 实现认证服务接口
type authService struct {
	userRepo repository.UserRepository // 用户数据访问接口
	jwtCfg   *config.JWTConfig         // 建议改为 jwtConfig
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, jwtCfg *config.JWTConfig) *authService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// Register 实现用户注册逻辑
func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) error {
	// 检查用户是否已存在
	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return errors.ErrUserExists
	}
	if err != errors.ErrUserNotFound {
		return err
	}

	// 创建新用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
	}
	// 设置加密密码
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

// Login 实现用户登录逻辑
func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (string, *auth.UserInfo, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return "", nil, errors.ErrInvalidCredentials
		}
		return "", nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return "", nil, errors.ErrInvalidCredentials
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, s.jwtCfg)
	if err != nil {
		return "", nil, err
	}

	// 构造用户信息
	userInfo := &auth.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),  // 格式化时间
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),  // 格式化时间
	}

	return token, userInfo, nil
}
