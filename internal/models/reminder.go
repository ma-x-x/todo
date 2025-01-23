package models

import (
	"errors"
	"time"
)

// ReminderType 提醒类型
type ReminderType int

const (
	ReminderTypeOnce   ReminderType = iota + 1 // 一次性提醒
	ReminderTypeDaily                          // 每日提醒
	ReminderTypeWeekly                         // 每周提醒
)

// NotifyType 通知类型
type NotifyType int

const (
	NotifyTypeEmail NotifyType = iota + 1 // 邮件通知
	NotifyTypePush                        // 推送通知
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
// @Description 待办事项提醒
type Reminder struct {
	Base
	TodoID     uint         `json:"todo_id" gorm:"not null"`
	RemindAt   time.Time    `json:"remind_at" gorm:"not null"`
	RemindType ReminderType `json:"remind_type" gorm:"type:int;not null"`
	NotifyType NotifyType   `json:"notify_type" gorm:"type:int;not null"`
	Status     bool         `json:"status" gorm:"default:false"` // false: 未提醒, true: 已提醒
}

func (r ReminderType) Int() int {
	return int(r)
}

func (n NotifyType) Int() int {
	return int(n)
}

// ParseReminderType 从字符串解析提醒类型
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

// ParseNotifyType 从字符串解析通知类型
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
