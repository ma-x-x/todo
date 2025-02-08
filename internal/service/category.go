package service

import (
	"context"
	"todo/api/v1/dto/category"
	"todo/internal/models"
)

// CategoryService 定义分类服务接口
type CategoryService interface {
	// Create 创建分类
	Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error)

	// List 获取用户的分类列表
	List(ctx context.Context, userID uint) ([]*models.Category, error)

	// Get 获取分类详情
	Get(ctx context.Context, id, userID uint) (*models.Category, error)

	// Update 更新分类信息
	Update(ctx context.Context, id, userID uint, req *category.UpdateRequest) error

	// Delete 删除分类
	Delete(ctx context.Context, id, userID uint) error

	// CreateCategory 创建新的分类
	CreateCategory(ctx context.Context, req *category.CreateRequest, userID uint) (*models.Category, error)
}
