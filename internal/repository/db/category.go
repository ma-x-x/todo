// Package db 提供数据库访问的具体实现
package db

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// categoryRepository 实现分类数据库操作的结构体
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储的实例
// db: 数据库连接实例
// 返回: 分类仓储实例
func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db: db}
}

// Create 在数据库中创建新的分类记录
// ctx: 上下文信息
// category: 要创建的分类信息
// 返回: error 创建过程中的错误信息
func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID 根据ID从数据库获取分类信息
// ctx: 上下文信息
// id: 分类ID
// 返回: (*models.Category, error) 分类信息和可能的错误
func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) ListByUserID(ctx context.Context, userID uint) ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}
