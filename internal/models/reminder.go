package models

import (
	"errors"
	"time"
)

// ReminderType 提醒类型
// 定义不同的提醒周期类型
type ReminderType int

const (
	ReminderTypeOnce   ReminderType = 1 // 一次性提醒
	ReminderTypeDaily  ReminderType = 2 // 每日提醒
	ReminderTypeWeekly ReminderType = 3 // 每周提醒
)

// NotifyType 通知类型
// 定义不同的通知方式
type NotifyType int

const (
	NotifyTypeEmail NotifyType = 1 // 邮件通知
	NotifyTypePush  NotifyType = 2 // 推送通知
)

// String 方法用于将 ReminderType 转换为字符串
func (r ReminderType) String() string {
	switch r {
	case ReminderTypeOnce:
		return "once"
	case ReminderTypeDaily:
		return "daily"
	case ReminderTypeWeekly:
		return "weekly"
	default:
		return "unknown"
	}
}

// String 方法用于将 NotifyType 转换为字符串
func (n NotifyType) String() string {
	switch n {
	case NotifyTypeEmail:
		return "email"
	case NotifyTypePush:
		return "push"
	default:
		return "unknown"
	}
}

// Reminder 提醒模型
// 存储待办事项的提醒信息，包括提醒时间、提醒类型和通知方式等
type Reminder struct {
	Base
	TodoID     uint         `json:"todo_id" gorm:"not null"`              // 关联的待办事项ID
	RemindAt   time.Time    `json:"remind_at" gorm:"not null"`            // 提醒时间
	RemindType ReminderType `json:"remind_type" gorm:"type:int;not null"` // 提醒类型
	NotifyType NotifyType   `json:"notify_type" gorm:"type:int;not null"` // 通知类型
	Status     bool         `json:"status" gorm:"default:false"`          // 提醒状态，false表示未提醒，true表示已提醒
}

func (r ReminderType) Int() int {
	return int(r)
}

func (n NotifyType) Int() int {
	return int(n)
}

// ParseReminderType 将字符串转换为提醒类型
func ParseReminderType(s string) (ReminderType, error) {
	switch s {
	case "once":
		return ReminderTypeOnce, nil
	case "daily":
		return ReminderTypeDaily, nil
	case "weekly":
		return ReminderTypeWeekly, nil
	default:
		return 0, errors.New("invalid reminder type")
	}
}

// ParseNotifyType 将字符串转换为通知类型
func ParseNotifyType(s string) (NotifyType, error) {
	switch s {
	case "email":
		return NotifyTypeEmail, nil
	case "push":
		return NotifyTypePush, nil
	default:
		return 0, errors.New("invalid notify type")
	}
}
