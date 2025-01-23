package impl

import (
	"context"
	"todo-demo/api/v1/dto/reminder"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/errors"
)

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
	// 验证待办事项是否存在且属于当前用户
	todo, err := s.todoRepo.GetByID(ctx, req.TodoID)
	if err != nil {
		return 0, err
	}
	if todo.UserID != userID {
		return 0, errors.ErrForbidden
	}

	// 将请求中的字符串类型转换为枚举类型
	reminderType, err := models.ParseReminderType(req.RemindType)
	if err != nil {
		return 0, err
	}

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

func (s *ReminderService) GetByID(ctx context.Context, userID, reminderID uint) (*models.Reminder, error) {
	reminder, err := s.reminderRepo.GetByID(ctx, reminderID)
	if err != nil {
		return nil, err
	}

	// 验证提醒是否属于当前用户
	todo, err := s.todoRepo.GetByID(ctx, reminder.TodoID)
	if err != nil {
		return nil, err
	}
	if todo.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return reminder, nil
}

func (s *ReminderService) ListByTodoID(ctx context.Context, userID, todoID uint) ([]*models.Reminder, error) {
	// 验证待办事项是否属于当前用户
	todo, err := s.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if todo.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return s.reminderRepo.ListByTodoID(ctx, todoID)
}

func (s *ReminderService) GetReminderRepo() repository.ReminderRepository {
	return s.reminderRepo
}

func (s *ReminderService) Update(ctx context.Context, reminder *models.Reminder) error {
	return s.reminderRepo.Update(ctx, reminder)
}

func (s *ReminderService) Delete(ctx context.Context, id uint) error {
	return s.reminderRepo.Delete(ctx, id)
}
