package service

import (
	"context"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/models"
)

// TodoService 定义待办事项相关的业务接口
type TodoService interface {
	// Create 创建待办事项
	Create(ctx context.Context, userID uint, req *todo.CreateRequest) (uint, error)

	// List 获取用户的待办事项列表
	List(ctx context.Context, userID uint) ([]*models.Todo, error)

	// Get 获取单个待办事项详情
	Get(ctx context.Context, id, userID uint) (*models.Todo, error)

	// Update 更新待办事项
	Update(ctx context.Context, id, userID uint, req *todo.UpdateRequest) error

	// Delete 删除待办事项
	Delete(ctx context.Context, id, userID uint) error
}
