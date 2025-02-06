package reminder

import "todo/internal/models"

// ListResponse 提醒列表响应
type ListResponse struct {
	Items []*models.Reminder `json:"items"`
	Total int64             `json:"total"`
}
