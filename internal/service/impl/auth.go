package impl

import (
	"context"
	"todo-demo/api/v1/dto/auth"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/config"
	"todo-demo/pkg/errors"
	"todo-demo/pkg/utils"
)

// authService 实现认证服务接口
type authService struct {
	userRepo repository.UserRepository // 用户数据访问接口
	jwtCfg   *config.JWTConfig         // JWT配置信息
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
func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (string, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return "", errors.ErrInvalidCredentials
		}
		return "", err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return "", errors.ErrInvalidCredentials
	}

	// 生成JWT令牌
	return utils.GenerateToken(user.ID, s.jwtCfg)
}
