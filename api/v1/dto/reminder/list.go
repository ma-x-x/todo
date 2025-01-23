package reminder

import "todo-demo/internal/models"

// ListResponse 提醒列表响应
type ListResponse struct {
	Total     int64              `json:"total"`
	Reminders []*models.Reminder `json:"reminders"`
}
