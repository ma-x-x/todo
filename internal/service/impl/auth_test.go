package impl

import (
	"context"
	"testing"
	"todo/api/v1/dto/auth"
	"todo/internal/models"
	"todo/pkg/config"
	"todo/pkg/errors"
)

// mockUserRepo 模拟用户仓储接口，用于单元测试
type mockUserRepo struct {
	users map[string]*models.User // 用户数据的内存存储
}

// newMockUserRepo 创建一个新的模拟用户仓储实例
func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*models.User),
	}
}

// Create 实现创建用户的模拟方法
func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	// 检查用户是否已存在
	if _, exists := m.users[user.Username]; exists {
		return errors.ErrUserExists
	}
	m.users[user.Username] = user
	return nil
}

// GetByUsername 实现根据用户名查询用户的模拟方法
func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.ErrUserNotFound
}

func (m *mockUserRepo) Delete(ctx context.Context, id uint) error {
	for _, user := range m.users {
		if user.ID == id {
			delete(m.users, user.Username)
			return nil
		}
	}
	return errors.ErrUserNotFound
}

func (m *mockUserRepo) Update(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return errors.ErrUserNotFound
	}
	m.users[user.Username] = user
	return nil
}

// TestAuthService_Register 测试用户注册功能
func TestAuthService_Register(t *testing.T) {
	// 初始化测试环境
	userRepo := newMockUserRepo()
	jwtCfg := &config.JWTConfig{
		Secret:      "test_secret", // 测试用的JWT密钥
		ExpireHours: 24,            // token过期时间
		Issuer:      "test",        // 令牌签发者
	}
	authService := NewAuthService(userRepo, jwtCfg)

	// 定义测试用例
	tests := []struct {
		name    string                // 测试用例名称
		req     *auth.RegisterRequest // 注册请求
		wantErr error                 // 期望的错误
	}{
		{
			name: "success", // 成功注册的用例
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "duplicate_username", // 重复用户名的用例
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test2@example.com",
			},
			wantErr: errors.ErrUserExists,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Register(context.Background(), tt.req)
			if err != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
