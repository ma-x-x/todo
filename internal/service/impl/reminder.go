package impl

import (
	"context"
	"time"
	"todo/api/v1/dto/reminder"
	"todo/internal/models"
	"todo/internal/repository/interfaces"
	"todo/pkg/errors"
)

// ReminderService 提醒服务实现
type ReminderService struct {
	reminderRepo interfaces.ReminderRepository
	todoRepo     interfaces.TodoRepository
}

func NewReminderService(reminderRepo interfaces.ReminderRepository, todoRepo interfaces.TodoRepository) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
		todoRepo:     todoRepo,
	}
}

func (s *ReminderService) Create(ctx context.Context, userID uint, req *reminder.CreateRequest) (uint, error) {
	todo, err := s.todoRepo.GetByID(ctx, req.TodoID)
	if err != nil {
		return 0, err
	}
	if todo.UserID != userID {
		return 0, errors.ErrNoPermission
	}

	remindAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.RemindAt)
	if err != nil {
		return 0, errors.New("无效的时间格式")
	}

	remindType, err := models.ParseRemindType(req.RemindType)
	if err != nil {
		return 0, err
	}

	notifyType, err := models.ParseNotifyType(req.NotifyType)
	if err != nil {
		return 0, err
	}

	// 创建提醒
	reminder := &models.Reminder{
		TodoID:     req.TodoID,
		RemindAt:   remindAt,
		RemindType: remindType,
		NotifyType: notifyType,
		Status:     false,
	}

	// 验证提醒数据
	if err := reminder.Validate(); err != nil {
		return 0, err
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

func (s *ReminderService) GetReminderRepo() interfaces.ReminderRepository {
	return s.reminderRepo
}

func (s *ReminderService) Update(ctx context.Context, id, userID uint, req *reminder.UpdateRequest) error {
	reminder, err := s.reminderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	todo, err := s.todoRepo.GetByID(ctx, reminder.TodoID)
	if err != nil {
		return err
	}
	if todo.UserID != userID {
		return errors.ErrNoPermission
	}

	// 更新提醒信息
	if req.RemindAt != "" {
		remindAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.RemindAt)
		if err != nil {
			return errors.New("无效的时间格式")
		}
		reminder.RemindAt = remindAt
	}
	if req.RemindType != "" {
		remindType, err := models.ParseRemindType(req.RemindType)
		if err != nil {
			return err
		}
		reminder.RemindType = remindType
	}
	if req.NotifyType != "" {
		notifyType, err := models.ParseNotifyType(req.NotifyType)
		if err != nil {
			return err
		}
		reminder.NotifyType = notifyType
	}

	// 验证提醒数据
	if err := reminder.Validate(); err != nil {
		return err
	}

	return s.reminderRepo.Update(ctx, reminder)
}

func (s *ReminderService) Delete(ctx context.Context, id, userID uint) error {
	reminder, err := s.Get(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.reminderRepo.Delete(ctx, reminder.ID)
}
