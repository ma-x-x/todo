package repository

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id uint) (*models.Todo, error)
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*models.Todo, int64, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id uint) error
}

type todoRepo struct {
	db *gorm.DB
}

func (r *todoRepo) Create(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepo) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrTodoNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepo) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*models.Todo, int64, error) {
	var todos []*models.Todo
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Todo{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (r *todoRepo) Update(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Todo{}, id).Error
}
