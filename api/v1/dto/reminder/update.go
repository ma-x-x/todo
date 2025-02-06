package reminder

import (
	"errors"
	"time"
)

// UpdateRequest 更新提醒请求
type UpdateRequest struct {
	RemindAt   time.Time `json:"remindAt" binding:"required"`
	RemindType string    `json:"remindType" binding:"required,oneof=once daily weekly"`
	NotifyType string    `json:"notifyType" binding:"required,oneof=email push"`
}

// Validate 验证请求参数
func (r *UpdateRequest) Validate() error {
	if r.RemindAt.Before(time.Now()) {
		return errors.New("提醒时间不能是过去时间")
	}
	return nil
}

// UpdateResponse 更新提醒响应
type UpdateResponse struct {
	Message string `json:"message"`
}
