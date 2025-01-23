package service

import (
	"context"
	"todo-demo/internal/models"
)

// TodoService 定义待办事项相关的业务接口
type TodoService interface {
	// Create 创建待办事项
	// ctx: 上下文信息
	// userID: 创建者ID
	// title: 标题
	// description: 描述
	// priority: 优先级
	// categoryID: 分类ID(可选)
	Create(ctx context.Context, userID uint, title, description string,
		priority models.Priority, categoryID *uint) (*models.Todo, error)

	// List 获取用户的待办事项列表
	List(ctx context.Context, userID uint) ([]*models.Todo, error)

	// Get 获取单个待办事项详情
	Get(ctx context.Context, id, userID uint) (*models.Todo, error)

	// Update 更新待办事项
	// req: 更新请求，包含要更新的字段
	Update(ctx context.Context, id, userID uint, req interface{}) error

	// Delete 删除待办事项
	Delete(ctx context.Context, id, userID uint) error
}
