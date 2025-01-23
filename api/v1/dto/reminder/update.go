package reminder

import (
	"errors"
	"time"
)

// UpdateRequest 更新提醒请求
type UpdateRequest struct {
	// RemindAt 提醒时间
	RemindAt time.Time `json:"remind_at" binding:"required"`
	// RemindType 提醒类型：once(单次)/daily(每日)/weekly(每周)
	RemindType string `json:"remind_type" binding:"required,oneof=once daily weekly"`
	// NotifyType 通知方式：email(邮件)/push(推送)
	NotifyType string `json:"notify_type" binding:"required,oneof=email push"`
}

// Validate 验证请求参数
func (r *UpdateRequest) Validate() error {
	// 验证提醒时间不能为过去时间
	if r.RemindAt.Before(time.Now()) {
		return errors.New("提醒时间不能是过去时间")
	}
	return nil
}

// UpdateResponse 更新提醒响应
type UpdateResponse struct {
	// Message 响应消息
	Message string `json:"message"`
}
