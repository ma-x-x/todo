package service

import (
	"context"
	"time"
	"todo-demo/internal/models"
)

// ReminderService 定义提醒相关的业务接口
type ReminderService interface {
	// Create 创建提醒
	// remindType: 提醒类型(一次性/每日/每周)
	// notifyType: 通知方式(邮件/推送)
	Create(ctx context.Context, userID, todoID uint, remindAt time.Time,
		remindType models.ReminderType, notifyType models.NotifyType) (*models.Reminder, error)

	// List 获取待办事项的提醒列表
	List(ctx context.Context, todoID uint) ([]*models.Reminder, error)

	// Get 获取提醒详情
	Get(ctx context.Context, id uint) (*models.Reminder, error)

	// Update 更新提醒信息
	Update(ctx context.Context, id uint, remindAt time.Time,
		remindType models.ReminderType, notifyType models.NotifyType) error

	// Delete 删除提醒
	Delete(ctx context.Context, id uint) error
}
