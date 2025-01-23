package impl

import (
	"context"
	"todo-demo/api/v1/dto/category"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/errors"
)

// CategoryService 分类服务结构体
// 负责处理所有与分类相关的业务逻辑
type CategoryService struct {
	categoryRepo repository.CategoryRepository // 分类数据仓库接口
}

// NewCategoryService 创建一个新的分类服务实例
// @param categoryRepo - 分类仓库实现
// @return 返回分类服务实例
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// Create 创建新的分类
// @param ctx - 上下文信息
// @param userID - 用户ID
// @param req - 创建分类的请求数据
// @return 返回新创建的分类ID和可能的错误
func (s *CategoryService) Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error) {
	category := &models.Category{
		Name:   req.Name,  // 分类名称
		Color:  req.Color, // 分类颜色
		UserID: userID,    // 所属用户ID
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return 0, err
	}

	return category.ID, nil
}

// GetByID 根据ID获取分类信息
// @param ctx - 上下文信息
// @param userID - 用户ID
// @param categoryID - 分类ID
// @return 返回分类信息和可能的错误
func (s *CategoryService) GetByID(ctx context.Context, userID, categoryID uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// 验证分类是否属于当前用户
	if category.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return category, nil
}

// List 获取用户的所有分类列表
// @param ctx - 上下文信息
// @param userID - 用户ID
// @return 返回分类列表和可能的错误
func (s *CategoryService) List(ctx context.Context, userID uint) ([]*models.Category, error) {
	return s.categoryRepo.ListByUserID(ctx, userID)
}

// Update 更新分类信息
// @param ctx - 上下文信息
// @param userID - 用户ID
// @param categoryID - 分类ID
// @param req - 更新分类的请求数据
// @return 返回可能的错误
func (s *CategoryService) Update(ctx context.Context, userID, categoryID uint, req *category.UpdateRequest) error {
	category, err := s.GetByID(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	// 只更新提供的字段
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Color != nil {
		category.Color = *req.Color
	}

	return s.categoryRepo.Update(ctx, category)
}

// Delete 删除分类
// @param ctx - 上下文信息
// @param userID - 用户ID
// @param categoryID - 分类ID
// @return 返回可能的错误
func (s *CategoryService) Delete(ctx context.Context, userID, categoryID uint) error {
	category, err := s.GetByID(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, category.ID)
}
