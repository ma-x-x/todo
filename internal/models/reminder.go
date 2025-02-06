package models

import (
	"errors"
	"time"
)

// 提醒类型常量
const (
	RemindTypeOnceStr   = "once"
	RemindTypeDailyStr  = "daily"
	RemindTypeWeeklyStr = "weekly"
)

// 通知类型常量
const (
	NotifyTypeEmailStr = "email"
	NotifyTypePushStr  = "push"
)

// Reminder 提醒模型
// 存储待办事项的提醒信息，包括提醒时间、提醒类型和通知方式等
type Reminder struct {
	Base
	TodoID     uint      `json:"todoId" gorm:"column:todo_id;not null;type:bigint unsigned;index"` // 关联的待办事项ID
	RemindAt   time.Time `json:"remindAt" gorm:"column:remind_at;not null;type:datetime"`         // 提醒时间
	RemindType string    `json:"remindType" gorm:"column:remind_type;not null;type:varchar(10)"`  // 提醒类型
	NotifyType string    `json:"notifyType" gorm:"column:notify_type;not null;type:varchar(10)"`  // 通知类型
	Status     bool      `json:"status" gorm:"column:status;not null;default:false"`              // 提醒状态
	Todo       *Todo     `json:"todo" gorm:"foreignKey:TodoID"`                                   // 关联的待办事项
}

// TableName 指定表名
func (Reminder) TableName() string {
	return "reminders"
}

// Validate 验证提醒数据
func (r *Reminder) Validate() error {
	// 验证 RemindType
	switch r.RemindType {
	case RemindTypeOnceStr, RemindTypeDailyStr, RemindTypeWeeklyStr:
	default:
		return errors.New("无效的提醒类型")
	}

	// 验证 NotifyType
	switch r.NotifyType {
	case NotifyTypeEmailStr, NotifyTypePushStr:
	default:
		return errors.New("无效的通知类型")
	}

	return nil
}
