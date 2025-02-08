package interfaces

import (
	"context"
	"todo/internal/models"
)

// ReminderRepository 提醒仓储接口
type ReminderRepository interface {
	// Create 创建提醒
	Create(ctx context.Context, reminder *models.Reminder) error

	// GetByID 根据ID获取提醒
	GetByID(ctx context.Context, id uint) (*models.Reminder, error)

	// ListByTodoID 获取待办事项的所有提醒
	ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error)

	// Update 更新提醒
	Update(ctx context.Context, reminder *models.Reminder) error

	// Delete 删除提醒
	Delete(ctx context.Context, id uint) error

	// GetPendingReminders 获取待处理的提醒
	GetPendingReminders(ctx context.Context) ([]*models.Reminder, error)
}
