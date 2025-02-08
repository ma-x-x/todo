package reminder

import (
	"time"
	"todo/internal/models"
)

// ReminderResponse 提醒响应
type ReminderResponse struct {
	ID         uint         `json:"id"`
	TodoID     uint         `json:"todo_id"`
	RemindAt   time.Time    `json:"remind_at"`
	RemindType string       `json:"remind_type"`
	NotifyType string       `json:"notify_type"`
	Status     bool         `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Todo       *models.Todo `json:"todo,omitempty"`
}

// ListResponse 提醒列表响应
type ListResponse struct {
	Total int                 `json:"total"`
	Items []*ReminderResponse `json:"items"`
}

// ConvertToResponse 将模型转换为响应
func ConvertToResponse(reminder *models.Reminder) *ReminderResponse {
	if reminder == nil {
		return nil
	}

	return &ReminderResponse{
		ID:         reminder.ID,
		TodoID:     reminder.TodoID,
		RemindAt:   reminder.RemindAt,
		RemindType: reminder.RemindType.String(),
		NotifyType: reminder.NotifyType.String(),
		Status:     reminder.Status,
		CreatedAt:  reminder.CreatedAt,
		UpdatedAt:  reminder.UpdatedAt,
		Todo:       reminder.Todo,
	}
}
