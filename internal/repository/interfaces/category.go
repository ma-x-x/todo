package interfaces

import (
	"context"
	"todo/internal/models"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, category *models.Category) error

	// GetByID 根据ID获取分类
	GetByID(ctx context.Context, id uint) (*models.Category, error)

	// GetByIDAndUserID 根据ID和用户ID获取分类
	GetByIDAndUserID(ctx context.Context, id, userID uint) (*models.Category, error)

	// ListByUserID 获取用户的所有分类
	ListByUserID(ctx context.Context, userID uint) ([]*models.Category, error)

	// Update 更新分类
	Update(ctx context.Context, category *models.Category) error

	// Delete 删除分类
	Delete(ctx context.Context, id uint) error
}
