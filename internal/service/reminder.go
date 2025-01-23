package service

import (
	"context"
	"time"
	"todo-demo/internal/models"
)

type ReminderService interface {
	Create(ctx context.Context, userID, todoID uint, remindAt time.Time,
		remindType models.ReminderType, notifyType models.NotifyType) (*models.Reminder, error)
	List(ctx context.Context, todoID uint) ([]*models.Reminder, error)
	Get(ctx context.Context, id uint) (*models.Reminder, error)
	Update(ctx context.Context, id uint, remindAt time.Time,
		remindType models.ReminderType, notifyType models.NotifyType) error
	Delete(ctx context.Context, id uint) error
}
