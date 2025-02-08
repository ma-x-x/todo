package impl

import (
	"context"
	"time"
	"todo/internal/models"
	"todo/internal/repository/interfaces"

	"gorm.io/gorm"
)

// ReminderRepository 提醒仓储实现
type ReminderRepository struct {
	*BaseRepository
}

// NewReminderRepository 创建提醒仓储实例
func NewReminderRepository(db *gorm.DB) interfaces.ReminderRepository {
	return &ReminderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建提醒
func (r *ReminderRepository) Create(ctx context.Context, reminder *models.Reminder) error {
	return r.BaseRepository.Create(ctx, reminder)
}

// GetByID 根据ID获取提醒
func (r *ReminderRepository) GetByID(ctx context.Context, id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.WithContext(ctx).First(&reminder, id).Error
	return &reminder, r.handleError(err, "reminder")
}

// ListByTodoID 获取待办事项的所有提醒
func (r *ReminderRepository) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	if err := r.List(ctx, 0, -1, &reminders, "todo_id = ?", todoID); err != nil {
		return nil, err
	}
	return reminders, nil
}

// Update 更新提醒
func (r *ReminderRepository) Update(ctx context.Context, reminder *models.Reminder) error {
	return r.BaseRepository.Update(ctx, reminder)
}

// Delete 删除提醒
func (r *ReminderRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &models.Reminder{ID: id})
}

// GetPendingReminders 获取待处理的提醒
func (r *ReminderRepository) GetPendingReminders(ctx context.Context) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	now := time.Now()
	if err := r.List(ctx, 0, -1, &reminders, "status = ? AND remind_at <= ?", false, now); err != nil {
		return nil, err
	}
	return reminders, nil
}
