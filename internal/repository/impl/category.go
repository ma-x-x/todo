package impl

import (
	"context"
	"todo/internal/models"
	"todo/internal/repository/interfaces"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository 创建分类仓储实例
func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.BaseRepository.Create(ctx, category)
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	return &category, r.handleError(err, "category")
}

// GetByIDAndUserID 根据ID和用户ID获取分类
func (r *CategoryRepository) GetByIDAndUserID(ctx context.Context, id, userID uint) (*models.Category, error) {
	var category models.Category
	if err := r.GetDB(ctx).Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// ListByUserID 获取用户的所有分类
func (r *CategoryRepository) ListByUserID(ctx context.Context, userID uint) ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.List(ctx, 0, -1, &categories, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类
func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.BaseRepository.Update(ctx, category)
}

// Delete 删除分类
func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &models.Category{Base: models.Base{ID: id}})
}
