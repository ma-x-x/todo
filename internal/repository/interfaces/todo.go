package interfaces

import (
	"context"
	"todo/internal/models"
)

// TodoRepository 待办事项仓储接口
type TodoRepository interface {
	// Create 创建待办事项
	Create(ctx context.Context, todo *models.Todo) error

	// GetByID 根据ID获取待办事项
	GetByID(ctx context.Context, id uint) (*models.Todo, error)

	// GetByIDAndUserID 根据ID和用户ID获取待办事项
	GetByIDAndUserID(ctx context.Context, id, userID uint) (*models.Todo, error)

	// ListByUserID 获取用户的所有待办事项
	ListByUserID(ctx context.Context, userID uint) ([]*models.Todo, error)

	// Update 更新待办事项
	Update(ctx context.Context, todo *models.Todo) error

	// Delete 删除待办事项
	Delete(ctx context.Context, id uint) error
}
