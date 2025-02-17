package reminder

import "todo/internal/models"

// ListResponse 提醒列表响应
type ListResponse struct {
	Total int                `json:"total"`
	Items []ReminderResponse `json:"items"`
}

// ConvertToResponse 将模型转换为响应
func ConvertToResponse(reminder *models.Reminder) *ReminderResponse {
	if reminder == nil {
		return nil
	}

	return &ReminderResponse{
		ID:         reminder.ID,
		TodoID:     reminder.TodoID,
		RemindAt:   reminder.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
		RemindType: reminder.RemindType.String(),
		NotifyType: reminder.NotifyType.String(),
		Status:     reminder.Status,
		CreatedAt:  reminder.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  reminder.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
