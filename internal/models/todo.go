package models

// Priority 优先级
type Priority int

const (
	PriorityLow    Priority = 1 // 低优先级
	PriorityMedium Priority = 2 // 中优先级
	PriorityHigh   Priority = 3 // 高优先级
)

// Todo 待办事项模型
// @Description 待办事项
type Todo struct {
	Base
	Title       string     `json:"title" gorm:"size:128;not null"`
	Description string     `json:"description" gorm:"size:1024"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	Priority    Priority   `json:"priority" gorm:"default:2"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        User       `gorm:"foreignKey:UserID" json:"-"`
	CategoryID  *uint      `json:"category_id" gorm:"default:null"`
	Category    *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Reminders   []Reminder `json:"reminders,omitempty" gorm:"foreignKey:TodoID"`
}
