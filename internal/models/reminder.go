package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
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
	gorm.Model       `json:"-"`                                              // 嵌入 gorm.Model 但在 JSON 中隐藏
	ID        uint   `json:"id" gorm:"primarykey"`                          // 主键ID
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`         // 创建时间
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`         // 更新时间
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"` // 软删除时间
	TodoID    uint      `json:"todoId" gorm:"column:todo_id;not null;type:bigint unsigned;index:idx_reminders_todo_id"`    // 关联的待办事项ID
	RemindAt  time.Time `json:"remindAt" gorm:"column:remind_at;not null;type:datetime;index:idx_reminders_time"`          // 提醒时间
	RemindType string   `json:"remindType" gorm:"column:remind_type;not null;type:varchar(10)"`   // 提醒类型
	NotifyType string   `json:"notifyType" gorm:"column:notify_type;not null;type:varchar(10)"`   // 通知类型
	Status     bool     `json:"status" gorm:"column:status;not null;default:false"`               // 提醒状态
	Todo       *Todo    `json:"todo,omitempty" gorm:"foreignKey:TodoID"`                         // 关联的待办事项
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
