// Package repository 实现数据访问层
package repository

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// ReminderRepository 定义提醒事项仓储接口
type ReminderRepository interface {
	// Create 创建新的提醒事项
	// ctx: 上下文信息
	// reminder: 提醒事项信息
	// 返回: error 创建过程中的错误信息
	Create(ctx context.Context, reminder *models.Reminder) error

	// GetByID 根据ID获取提醒事项
	// ctx: 上下文信息
	// id: 提醒事项ID
	// 返回: (*models.Reminder, error) 提醒事项信息和可能的错误
	GetByID(ctx context.Context, id uint) (*models.Reminder, error)

	// ListByTodoID 获取指定待办事项的所有提醒
	// ctx: 上下文信息
	// todoID: 待办事项ID
	// 返回: ([]*models.Reminder, error) 提醒事项列表和可能的错误
	ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error)

	// Update 更新提醒事项
	// ctx: 上下文信息
	// reminder: 需要更新的提醒事项信息
	// 返回: error 更新过程中的错误信息
	Update(ctx context.Context, reminder *models.Reminder) error

	// Delete 删除提醒事项
	// ctx: 上下文信息
	// id: 要删除的提醒事项ID
	// 返回: error 删除过程中的错误信息
	Delete(ctx context.Context, id uint) error
}

type reminderRepo struct {
	db *gorm.DB
}

func (r *reminderRepo) Create(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *reminderRepo) GetByID(ctx context.Context, id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := r.db.WithContext(ctx).First(&reminder, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrReminderNotFound
		}
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepo) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	if err := r.db.WithContext(ctx).Where("todo_id = ?", todoID).Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *reminderRepo) Update(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

func (r *reminderRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Reminder{}, id).Error
}
