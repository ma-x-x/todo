package reminder

import (
	"errors"
	"time"
	"todo-demo/internal/models"
)

// CreateRequest 创建提醒请求
type CreateRequest struct {
	// TodoID 关联的待办事项ID
	TodoID uint `json:"todo_id" binding:"required"`
	// RemindAt 提醒时间
	RemindAt time.Time `json:"remind_at" binding:"required"`
	// RemindType 提醒类型：once(单次)/daily(每日)/weekly(每周)
	RemindType string `json:"remind_type" binding:"required,oneof=once daily weekly"`
	// NotifyType 通知方式：email(邮件)/push(推送)
	NotifyType string `json:"notify_type" binding:"required,oneof=email push"`
}

// Validate 验证请求参数
func (r *CreateRequest) Validate() error {
	// 验证提醒时间不能为过去时间
	if r.RemindAt.Before(time.Now()) {
		return errors.New("提醒时间不能是过去时间")
	}

	// 验证待办事项ID
	if r.TodoID == 0 {
		return errors.New("待办事项ID不能为空")
	}

	return nil
}

// CreateResponse 创建提醒响应
type CreateResponse struct {
	// Reminder 创建成功后返回的提醒信息
	Reminder *models.Reminder `json:"reminder"`
}
