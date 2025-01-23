package impl

import (
	"context"
	"testing"
	"todo-demo/api/v1/dto/auth"
	"todo-demo/internal/models"
	"todo-demo/pkg/config"
	"todo-demo/pkg/errors"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.ErrUserExists
	}
	m.users[user.Username] = user
	return nil
}

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

func TestAuthService_Register(t *testing.T) {
	userRepo := newMockUserRepo()
	jwtCfg := &config.JWTConfig{
		Secret:      "test_secret",
		ExpireHours: 24,
		Issuer:      "test",
	}
	authService := NewAuthService(userRepo, jwtCfg)

	tests := []struct {
		name    string
		req     *auth.RegisterRequest
		wantErr error
	}{
		{
			name: "success",
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "duplicate_username",
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test2@example.com",
			},
			wantErr: errors.ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Register(context.Background(), tt.req)
			if err != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
