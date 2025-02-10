package impl

import (
	"context"
	"errors"
	"fmt"
	"todo/api/v1/dto/category"
	"todo/internal/models"
	"todo/internal/repository/interfaces"

	"gorm.io/gorm"
)

// CategoryService 分类服务实现
type CategoryService struct {
	categoryRepo interfaces.CategoryRepository
}

// NewCategoryService 创建一个新的分类服务实例
//
// Parameters:
//   - repo: 分类仓库实现
//
// Returns:
//   - *CategoryService: 返回分类服务实例
func NewCategoryService(repo interfaces.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: repo}
}

// Create 创建新的分类
//
// Parameters:
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - req: 创建分类的请求数据
//
// Returns:
//   - uint: 返回新创建的分类ID
//   - error: 可能的错误信息
func (s *CategoryService) Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		UserID:      userID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return 0, err
	}

	return category.ID, nil
}

// Get 根据ID获取分类信息
//
// Parameters:
//   - ctx: 上下文信息
//   - id: 分类ID
//   - userID: 用户ID
//
// Returns:
//   - *models.Category: 返回分类信息
//   - error: 可能的错误信息
func (s *CategoryService) Get(ctx context.Context, id, userID uint) (*models.Category, error) {
	return s.categoryRepo.GetByIDAndUserID(ctx, id, userID)
}

// List 获取用户的所有分类列表
//
// Parameters:
//   - ctx: 上下文信息
//   - userID: 用户ID
//
// Returns:
//   - []*models.Category: 返回分类列表
//   - error: 可能的错误信息
func (s *CategoryService) List(ctx context.Context, userID uint) ([]*models.Category, error) {
	return s.categoryRepo.ListByUserID(ctx, userID)
}

// Update 更新分类信息
//
// Parameters:
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - categoryID: 分类ID
//   - req: 更新分类的请求数据
//
// Returns:
//   - error: 可能的错误信息
func (s *CategoryService) Update(ctx context.Context, userID, id uint, req *category.UpdateRequest) error {
	// 先获取分类
	category, err := s.categoryRepo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	// 更新非空字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}

	return s.categoryRepo.Update(ctx, category)
}

// Delete 删除分类
//
// Parameters:
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - categoryID: 分类ID
//
// Returns:
//   - error: 可能的错误信息
func (s *CategoryService) Delete(ctx context.Context, id, userID uint) error {
	// 先检查分类是否存在且属于该用户
	_, err := s.categoryRepo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("分类不存在或无权限访问")
		}
		return fmt.Errorf("查询分类失败: %w", err)
	}

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	return nil
}

// CreateCategory 创建新的分类
func (s *CategoryService) CreateCategory(ctx context.Context, req *category.CreateRequest, userID uint) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		UserID:      userID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
