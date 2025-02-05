package reminder

import "todo-demo/internal/models"

// ListResponse 提醒列表响应
type ListResponse struct {
	Total     int64              `json:"total"`
	Items     []*models.Reminder `json:"items"`
}
