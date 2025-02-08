package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// RemindType 提醒类型
type RemindType string

const (
	RemindTypeOnce   RemindType = "once"
	RemindTypeDaily  RemindType = "daily"
	RemindTypeWeekly RemindType = "weekly"
)

// NotifyType 通知类型
type NotifyType string

const (
	NotifyTypeEmail NotifyType = "email"
	NotifyTypePush  NotifyType = "push"
)

// Reminder 提醒模型
// @Description 待办事项提醒
type Reminder struct {
	Base
	// 主键ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 关联的待办事项ID
	TodoID uint `json:"todoId" gorm:"not null;index"`
	// 提醒时间
	RemindAt time.Time `json:"remindAt" gorm:"not null;index"`
	// 提醒类型
	RemindType RemindType `json:"remindType" gorm:"size:20;not null"`
	// 通知类型
	NotifyType NotifyType `json:"notifyType" gorm:"size:20;not null"`
	// 提醒状态
	Status bool `json:"status" gorm:"default:false;index"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 关联的待办事项
	Todo *Todo `json:"todo,omitempty" gorm:"foreignKey:TodoID"`
}

// TableName 指定表名
func (Reminder) TableName() string {
	return "reminders"
}

// Validate 验证提醒数据
func (r *Reminder) Validate() error {
	if r.TodoID == 0 {
		return errors.New("待办事项ID不能为空")
	}

	if r.RemindAt.Before(time.Now()) {
		return errors.New("提醒时间不能是过去时间")
	}

	switch r.RemindType {
	case RemindTypeOnce, RemindTypeDaily, RemindTypeWeekly:
	default:
		return errors.New("无效的提醒类型")
	}

	switch r.NotifyType {
	case NotifyTypeEmail, NotifyTypePush:
	default:
		return errors.New("无效的通知类型")
	}

	return nil
}

// BeforeCreate 创建前的钩子函数
func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return r.Base.BeforeCreate(tx)
}

// BeforeUpdate 更新前的钩子函数
func (r *Reminder) BeforeUpdate(tx *gorm.DB) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return r.Base.BeforeUpdate(tx)
}

// String 实现 Stringer 接口
func (rt RemindType) String() string {
	return string(rt)
}

// String 实现 Stringer 接口
func (nt NotifyType) String() string {
	return string(nt)
}

// ParseRemindType 从字符串解析提醒类型
func ParseRemindType(s string) (RemindType, error) {
	switch RemindType(s) {
	case RemindTypeOnce, RemindTypeDaily, RemindTypeWeekly:
		return RemindType(s), nil
	default:
		return "", errors.New("无效的提醒类型")
	}
}

// ParseNotifyType 从字符串解析通知类型
func ParseNotifyType(s string) (NotifyType, error) {
	switch NotifyType(s) {
	case NotifyTypeEmail, NotifyTypePush:
		return NotifyType(s), nil
	default:
		return "", errors.New("无效的通知类型")
	}
}
