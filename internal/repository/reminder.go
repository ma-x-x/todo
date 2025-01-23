package repository

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder *models.Reminder) error
	GetByID(ctx context.Context, id uint) (*models.Reminder, error)
	ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error)
	Update(ctx context.Context, reminder *models.Reminder) error
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