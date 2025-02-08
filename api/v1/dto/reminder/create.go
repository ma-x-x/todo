package reminder

import (
	"time"
	"todo/internal/models"
)

// CreateRequest 创建提醒请求
// @Description 创建提醒的请求参数
type CreateRequest struct {
	// 待办事项ID
	TodoID uint `json:"todoId" binding:"required" example:"1"`
	// 提醒时间
	RemindAt string `json:"remindAt" binding:"required" example:"2024-02-08T17:12:40+08:00"`
	// 提醒类型 (once/daily/weekly)
	RemindType string `json:"remindType" binding:"required" example:"once"`
	// 通知类型 (email/push)
	NotifyType string `json:"notifyType" binding:"required" example:"email"`
}

// CreateResponse 创建提醒响应
type CreateResponse struct {
	ID         uint             `json:"id"`
	TodoID     uint             `json:"todoId"`
	RemindAt   time.Time        `json:"remindAt"`
	RemindType string           `json:"remindType"`
	NotifyType string           `json:"notifyType"`
	CreatedAt  time.Time        `json:"createdAt"`
	Todo       *models.Todo     `json:"todo"`     // 关联的待办事项信息
	Reminder   *models.Reminder `json:"reminder"` // 提醒详细信息
}
