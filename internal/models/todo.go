package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Priority 优先级类型
// 用于定义待办事项的优先级别
type Priority string

const (
	PriorityLow    Priority = "low"    // 低优先级
	PriorityMedium Priority = "medium" // 中优先级
	PriorityHigh   Priority = "high"   // 高优先级
)

// TodoStatus 待办事项状态
type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusProgress  TodoStatus = "in_progress"
	TodoStatusCompleted TodoStatus = "completed"
)

// Todo 待办事项模型
// @Description 待办事项
// 存储待办事项的详细信息，包括标题、描述、完成状态、优先级等
// 通过外键关联用户和分类信息
type Todo struct {
	Base
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"size:100;not null;index"`
	Description string     `json:"description" gorm:"size:500"`
	Status      TodoStatus `json:"status" gorm:"size:20;not null;default:pending;index"`
	Priority    Priority   `json:"priority" gorm:"size:20;default:medium;index"`
	DueDate     time.Time  `json:"dueDate"`
	CategoryID  *uint      `json:"categoryId"`
	UserID      uint       `json:"userId" gorm:"not null;index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Completed   bool       `json:"completed" gorm:"default:false;index"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
	Category    *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Reminders   []Reminder `json:"reminders,omitempty" gorm:"foreignKey:TodoID"`
}

// Validate 验证待办事项数据
func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("标题不能为空")
	}

	if t.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	if t.DueDate.Before(time.Now()) {
		return errors.New("截止时间不能是过去时间")
	}

	// 验证优先级
	switch t.Priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
	default:
		return errors.New("无效的优先级")
	}

	// 验证状态
	switch t.Status {
	case TodoStatusPending, TodoStatusProgress, TodoStatusCompleted:
	default:
		return errors.New("无效的状态")
	}

	return nil
}

// BeforeCreate 创建前的钩子函数
func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return t.Base.BeforeCreate(tx)
}

// BeforeUpdate 更新前的钩子函数
func (t *Todo) BeforeUpdate(tx *gorm.DB) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return t.Base.BeforeUpdate(tx)
}

// ParseTodoStatus 从字符串解析待办事项状态
func ParseTodoStatus(s string) (TodoStatus, error) {
	switch TodoStatus(s) {
	case TodoStatusPending, TodoStatusProgress, TodoStatusCompleted:
		return TodoStatus(s), nil
	default:
		return "", errors.New("无效的待办事项状态")
	}
}

// String 实现 Stringer 接口
func (p Priority) String() string {
	return string(p)
}

// String 实现 Stringer 接口
func (s TodoStatus) String() string {
	return string(s)
}
