package reminder

// ReminderResponse 提醒响应
type ReminderResponse struct {
	// 提醒ID
	ID uint `json:"id" example:"1"`
	// 待办事项ID
	TodoID uint `json:"todoId" example:"1"`
	// 提醒时间
	RemindAt string `json:"remindAt" example:"2024-02-08T17:10:54+08:00"`
	// 提醒类型
	RemindType string `json:"remindType" example:"once"`
	// 通知类型
	NotifyType string `json:"notifyType" example:"email"`
	// 提醒状态
	Status bool `json:"status" example:"false"`
	// 创建时间
	CreatedAt string `json:"createdAt" example:"2024-02-08T17:10:54+08:00"`
	// 更新时间
	UpdatedAt string `json:"updatedAt" example:"2024-02-08T17:10:54+08:00"`
}
