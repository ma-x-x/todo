package reminder

import "time"

// CreateRequest 创建提醒请求
type CreateRequest struct {
	TodoID     uint      `json:"todoId" binding:"required"`
	RemindAt   time.Time `json:"remindAt" binding:"required"`
	RemindType string    `json:"remindType" binding:"required,oneof=once daily weekly"`
	NotifyType string    `json:"notifyType" binding:"required,oneof=email push"`
}

// CreateResponse 创建提醒响应
type CreateResponse struct {
	ID uint `json:"id"`
}
