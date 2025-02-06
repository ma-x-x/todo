package service

import (
	"context"
	"todo/api/v1/dto/reminder"
	"todo/internal/models"
)

// ReminderService 提醒服务接口
type ReminderService interface {
	// Create 创建提醒
	// remindType: 提醒类型(一次性/每日/每周)
	// notifyType: 通知方式(邮件/推送)
	Create(ctx context.Context, userID uint, req *reminder.CreateRequest) (uint, error)

	// List 获取待办事项的提醒列表
	ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error)

	// Get 获取提醒详情
	Get(ctx context.Context, id, userID uint) (*models.Reminder, error)

	// Update 更新提醒信息
	Update(ctx context.Context, id, userID uint, req *reminder.UpdateRequest) error

	// Delete 删除提醒
	Delete(ctx context.Context, id, userID uint) error
}
