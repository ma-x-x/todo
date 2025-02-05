package impl

import (
	"context"
	"todo-demo/api/v1/dto/reminder"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/errors"
)

// ReminderService 提醒服务实现
type ReminderService struct {
	reminderRepo repository.ReminderRepository
	todoRepo     repository.TodoRepository
}

func NewReminderService(reminderRepo repository.ReminderRepository, todoRepo repository.TodoRepository) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
		todoRepo:     todoRepo,
	}
}

func (s *ReminderService) Create(ctx context.Context, userID uint, req *reminder.CreateRequest) (uint, error) {
	// 转换提醒类型
	reminderType, err := models.ParseReminderType(req.RemindType)
	if err != nil {
		return 0, err
	}

	// 转换通知类型
	notifyType, err := models.ParseNotifyType(req.NotifyType)
	if err != nil {
		return 0, err
	}

	reminder := &models.Reminder{
		TodoID:     req.TodoID,
		RemindAt:   req.RemindAt,
		RemindType: reminderType,
		NotifyType: notifyType,
	}

	if err := s.reminderRepo.Create(ctx, reminder); err != nil {
		return 0, err
	}
	return reminder.ID, nil
}

func (s *ReminderService) Get(ctx context.Context, id, userID uint) (*models.Reminder, error) {
	reminder, err := s.reminderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// 验证所有权（通过关联的Todo）
	todo, err := s.todoRepo.GetByID(ctx, reminder.TodoID)
	if err != nil {
		return nil, err
	}
	if todo.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return reminder, nil
}

func (s *ReminderService) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	return s.reminderRepo.ListByTodoID(ctx, todoID)
}

func (s *ReminderService) GetReminderRepo() repository.ReminderRepository {
	return s.reminderRepo
}

func (s *ReminderService) Update(ctx context.Context, id, userID uint, req *reminder.UpdateRequest) error {
	reminder, err := s.Get(ctx, id, userID)
	if err != nil {
		return err
	}

	// 转换提醒类型
	reminderType, err := models.ParseReminderType(req.RemindType)
	if err != nil {
		return err
	}

	// 转换通知类型
	notifyType, err := models.ParseNotifyType(req.NotifyType)
	if err != nil {
		return err
	}

	reminder.RemindAt = req.RemindAt
	reminder.RemindType = reminderType
	reminder.NotifyType = notifyType

	return s.reminderRepo.Update(ctx, reminder)
}

func (s *ReminderService) Delete(ctx context.Context, id, userID uint) error {
	reminder, err := s.Get(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.reminderRepo.Delete(ctx, reminder.ID)
}
