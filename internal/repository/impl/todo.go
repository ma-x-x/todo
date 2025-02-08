// Package impl 提供仓储层接口的具体实现
package impl

import (
	"context"
	"todo/internal/models"
	"todo/internal/repository/interfaces"

	"gorm.io/gorm"
)

// TodoRepository 待办事项仓储实现
type TodoRepository struct {
	*BaseRepository
}

// NewTodoRepository 创建待办事项仓储实例
func NewTodoRepository(db *gorm.DB) interfaces.TodoRepository {
	return &TodoRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建待办事项
func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	return r.BaseRepository.Create(ctx, todo)
}

// GetByID 根据ID获取待办事项
func (r *TodoRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.GetDB(ctx).Preload("Category").First(&todo, id).Error
	return &todo, r.handleError(err, "todo")
}

// GetByIDAndUserID 根据ID和用户ID获取待办事项
func (r *TodoRepository) GetByIDAndUserID(ctx context.Context, id, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.GetDB(ctx).Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	return &todo, r.handleError(err, "todo")
}

// ListByUserID 获取用户的所有待办事项
func (r *TodoRepository) ListByUserID(ctx context.Context, userID uint) ([]*models.Todo, error) {
	var todos []*models.Todo
	if err := r.List(ctx, 0, -1, &todos, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return todos, nil
}

// Update 更新待办事项
func (r *TodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	return r.BaseRepository.Update(ctx, todo)
}

// Delete 删除待办事项
func (r *TodoRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &models.Todo{Base: models.Base{ID: id}})
}
