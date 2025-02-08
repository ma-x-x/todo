package reminder

import "todo/internal/models"

// ReminderResponse 提醒响应
type ReminderResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
	DeletedAt  *string   `json:"deletedAt,omitempty"`
	TodoID     uint      `json:"todoId"`
	RemindAt   string    `json:"remindAt"`
	RemindType string    `json:"remindType"`
	NotifyType string    `json:"notifyType"`
	Status     bool      `json:"status"`
	Todo       *Todo     `json:"todo,omitempty"`
}

// ListResponse 提醒列表响应
type ListResponse struct {
	Items []*ReminderResponse `json:"items"`
	Total int64              `json:"total"`
}

// Todo 待办事项简要信息
type Todo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// ConvertToResponse 将模型转换为响应
func ConvertToResponse(reminder *models.Reminder) *ReminderResponse {
	resp := &ReminderResponse{
		ID:         reminder.ID,
		CreatedAt:  reminder.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  reminder.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		TodoID:     reminder.TodoID,
		RemindAt:   reminder.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
		RemindType: reminder.RemindType,
		NotifyType: reminder.NotifyType,
		Status:     reminder.Status,
	}

	if reminder.DeletedAt.Valid {
		deletedAt := reminder.DeletedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.DeletedAt = &deletedAt
	}

	if reminder.Todo != nil {
		resp.Todo = &Todo{
			ID:          reminder.Todo.ID,
			Title:       reminder.Todo.Title,
			Description: reminder.Todo.Description,
		}
	}

	return resp
}
