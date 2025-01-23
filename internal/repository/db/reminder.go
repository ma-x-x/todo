package db

import (
	"context"
	"todo-demo/internal/models"
	"todo-demo/pkg/errors"

	"gorm.io/gorm"
)

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *reminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

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

func (r *reminderRepository) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	if err := r.db.WithContext(ctx).Where("todo_id = ?", todoID).Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *reminderRepository) Update(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

func (r *reminderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Reminder{}, id).Error
}
