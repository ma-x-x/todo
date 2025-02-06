package reminder

import (
	"time"
	"todo/internal/models"
)

// CreateRequest 创建提醒请求
type CreateRequest struct {
	TodoID     uint      `json:"todoId" binding:"required"`
	RemindAt   time.Time `json:"remindAt" binding:"required"`
	RemindType string    `json:"remindType" binding:"required,oneof=once daily weekly"`
	NotifyType string    `json:"notifyType" binding:"required,oneof=email push"`
}

// CreateResponse 创建提醒响应
type CreateResponse struct {
	ID         uint           `json:"id"`
	TodoID     uint           `json:"todoId"`
	RemindAt   time.Time      `json:"remindAt"`
	RemindType string         `json:"remindType"`
	NotifyType string         `json:"notifyType"`
	CreatedAt  time.Time      `json:"createdAt"`
	Todo       *models.Todo     `json:"todo"`      // 关联的待办事项信息
	Reminder   *models.Reminder `json:"reminder"`  // 提醒详细信息
}
