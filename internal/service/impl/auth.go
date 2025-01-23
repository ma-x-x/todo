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

type authService struct {
	userRepo repository.UserRepository
	jwtCfg   *config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtCfg *config.JWTConfig) *authService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

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
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (string, error) {
	// 获取用户
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

	// 生成 token
	return utils.GenerateToken(user.ID, s.jwtCfg)
}
