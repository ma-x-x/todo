// Package db 提供数据库访问的具体实现
package db

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

// reminderRepository 实现提醒事项数据库操作的结构体
type reminderRepository struct {
	db *gorm.DB
}

// NewReminderRepository 创建提醒事项仓储的实例
// db: 数据库连接实例
// 返回: 提醒事项仓储实例
func NewReminderRepository(db *gorm.DB) *reminderRepository {
	return &reminderRepository{db: db}
}

// Create 在数据库中创建新的提醒事项记录
// ctx: 上下文信息
// reminder: 要创建的提醒事项信息
// 返回: error 创建过程中的错误信息
func (r *reminderRepository) Create(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

// GetByID 根据ID从数据库获取提醒事项信息
// ctx: 上下文信息
// id: 提醒事项ID
// 返回: (*models.Reminder, error) 提醒事项信息和可能的错误
func (r *reminderRepository) GetByID(ctx context.Context, id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := r.db.WithContext(ctx).First(&reminder, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrReminderNotFound
		}
		return nil, err
	}
	return &reminder, nil
}

// ListByTodoID 获取指定待办事项的所有提醒记录
// ctx: 上下文信息
// todoID: 待办事项ID
// 返回: ([]*models.Reminder, error) 提醒事项列表和可能的错误
func (r *reminderRepository) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	if err := r.db.WithContext(ctx).Where("todo_id = ?", todoID).Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}

// Update 更新数据库中的提醒事项记录
// ctx: 上下文信息
// reminder: 需要更新的提醒事项信息
// 返回: error 更新过程中的错误信息
func (r *reminderRepository) Update(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

// Delete 从数据库中删除提醒事项记录
// ctx: 上下文信息
// id: 要删除的提醒事项ID
// 返回: error 删除过程中的错误信息
func (r *reminderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Reminder{}, id).Error
}
