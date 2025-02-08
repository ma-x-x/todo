package interfaces

import (
	"context"
	"todo/internal/models"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *models.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *models.User) error

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error

	// List 获取用户列表
	List(ctx context.Context, offset, limit int) ([]*models.User, error)

	// Count 获取用户总数
	Count(ctx context.Context) (int64, error)
}
