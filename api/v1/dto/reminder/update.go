package reminder

import (
	"errors"
	"time"
)

// UpdateRequest 更新提醒请求
// @Description 更新提醒的请求参数
type UpdateRequest struct {
	// 提醒时间
	RemindAt string `json:"remindAt,omitempty" example:"2024-02-08T17:12:40+08:00"`
	// 提醒类型 (once/daily/weekly)
	RemindType string `json:"remindType,omitempty" example:"once"`
	// 通知类型 (email/push)
	NotifyType string `json:"notifyType,omitempty" example:"email"`
	// 提醒状态
	Status bool `json:"status,omitempty" example:"false"`
}

// Validate 验证请求参数
func (r *UpdateRequest) Validate() error {
	if r.RemindAt != "" {
		remindAt, err := time.Parse("2006-01-02T15:04:05Z07:00", r.RemindAt)
		if err != nil {
			return errors.New("无效的时间格式")
		}
		if remindAt.Before(time.Now()) {
			return errors.New("提醒时间不能是过去时间")
		}
	}
	return nil
}

// UpdateResponse 更新提醒响应
// @Description 更新提醒的响应
type UpdateResponse struct {
	// 响应消息
	Message string `json:"message" example:"更新成功"`
}
