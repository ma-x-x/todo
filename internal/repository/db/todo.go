package db

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *todoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo
	if err := r.db.WithContext(ctx).Preload("Category").First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrTodoNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*models.Todo, int64, error) {
	var todos []*models.Todo
	var total int64

	offset := (page - 1) * pageSize
	db := r.db.WithContext(ctx).Model(&models.Todo{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("Category").Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Todo{}, id).Error
}
