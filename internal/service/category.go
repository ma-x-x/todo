package service

import (
	"context"
	"todo-demo/api/v1/dto/category"
	"todo-demo/internal/models"
)

// CategoryService 分类服务接口
type CategoryService interface {
	// Create 创建分类
	// name: 分类名称
	// color: 分类颜色
	Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error)

	// List 获取用户的分类列表
	List(ctx context.Context, userID uint) ([]*models.Category, error)

	// Get 获取分类详情
	Get(ctx context.Context, id, userID uint) (*models.Category, error)

	// Update 更新分类信息
	Update(ctx context.Context, id, userID uint, req *category.UpdateRequest) error

	// Delete 删除分类
	Delete(ctx context.Context, id, userID uint) error
}
