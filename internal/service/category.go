package service

import (
	"context"
	"todo-demo/internal/models"
)

// CategoryService 定义分类相关的业务接口
type CategoryService interface {
	// Create 创建分类
	// name: 分类名称
	// color: 分类颜色
	Create(ctx context.Context, userID uint, name, color string) (*models.Category, error)

	// List 获取用户的分类列表
	List(ctx context.Context, userID uint) ([]*models.Category, error)

	// Get 获取分类详情
	Get(ctx context.Context, id, userID uint) (*models.Category, error)

	// Update 更新分类信息
	Update(ctx context.Context, id, userID uint, name, color string) error

	// Delete 删除分类
	Delete(ctx context.Context, id, userID uint) error
}
