// Package repository 实现数据访问层
package repository

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// TodoRepository 待办事项仓库接口
type TodoRepository interface {
	// Create 创建新的待办事项
	// ctx: 上下文信息
	// todo: 待办事项信息
	// 返回: error 创建过程中的错误信息
	Create(ctx context.Context, todo *models.Todo) error

	// GetByID 根据ID获取待办事项
	// ctx: 上下文信息
	// id: 待办事项ID
	// 返回: (*models.Todo, error) 待办事项信息和可能的错误
	GetByID(ctx context.Context, id uint) (*models.Todo, error)

	// ListByUserID 获取用户的待办事项列表
	// ctx: 上下文信息
	// userID: 用户ID
	// page: 页码
	// pageSize: 每页数量
	// 返回: ([]*models.Todo, int64, error) 待办事项列表、总数和可能的错误
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*models.Todo, int64, error)

	// Update 更新待办事项
	// ctx: 上下文信息
	// todo: 需要更新的待办事项信息
	// 返回: error 更新过程中的错误信息
	Update(ctx context.Context, todo *models.Todo) error

	// Delete 删除待办事项
	// ctx: 上下文信息
	// id: 要删除的待办事项ID
	// 返回: error 删除过程中的错误信息
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
