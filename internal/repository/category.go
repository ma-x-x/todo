// Package repository 实现数据访问层
package repository

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

// CategoryRepository 定义分类仓储接口
type CategoryRepository interface {
	// Create 创建新的分类
	// ctx: 上下文信息
	// category: 分类信息
	// 返回: error 创建过程中的错误信息
	Create(ctx context.Context, category *models.Category) error

	// GetByID 根据ID获取分类信息
	// ctx: 上下文信息
	// id: 分类ID
	// 返回: (*models.Category, error) 分类信息和可能的错误
	GetByID(ctx context.Context, id uint) (*models.Category, error)

	// ListByUserID 获取用户的所有分类
	// ctx: 上下文信息
	// userID: 用户ID
	// 返回: ([]*models.Category, error) 分类列表和可能的错误
	ListByUserID(ctx context.Context, userID uint) ([]*models.Category, error)

	// Update 更新分类信息
	// ctx: 上下文信息
	// category: 需要更新的分类信息
	// 返回: error 更新过程中的错误信息
	Update(ctx context.Context, category *models.Category) error

	// Delete 删除分类
	// ctx: 上下文信息
	// id: 要删除的分类ID
	// 返回: error 删除过程中的错误信息
	Delete(ctx context.Context, id uint) error
}

// categoryRepo 实现 CategoryRepository 接口
type categoryRepo struct {
	db *gorm.DB
}

func (r *categoryRepo) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepo) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) ListByUserID(ctx context.Context, userID uint) ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}
