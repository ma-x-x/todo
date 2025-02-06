// Package db 提供数据库访问的具体实现
package db

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// todoRepository 实现待办事项数据库操作的结构体
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建待办事项仓储的实例
// db: 数据库连接实例
// 返回: 待办事项仓储实例
func NewTodoRepository(db *gorm.DB) *todoRepository {
	return &todoRepository{db: db}
}

// Create 在数据库中创建新的待办事项记录
// ctx: 上下文信息
// todo: 要创建的待办事项信息
// 返回: error 创建过程中的错误信息
func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

// GetByID 根据ID从数据库获取待办事项信息，同时加载关联的分类信息
// ctx: 上下文信息
// id: 待办事项ID
// 返回: (*models.Todo, error) 待办事项信息和可能的错误
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

// ListByUserID 获取用户的待办事项列表，支持分页
// ctx: 上下文信息
// userID: 用户ID
// page: 页码
// pageSize: 每页数量
// 返回: ([]*models.Todo, int64, error) 待办事项列表、总数和可能的错误
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

// Update 更新数据库中的待办事项记录
// ctx: 上下文信息
// todo: 需要更新的待办事项信息
// 返回: error 更新过程中的错误信息
func (r *todoRepository) Update(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

// Delete 从数据库中删除待办事项记录
// ctx: 上下文信息
// id: 要删除的待办事项ID
// 返回: error 删除过程中的错误信息
func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Todo{}, id).Error
}
