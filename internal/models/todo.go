package models

// Priority 优先级类型
// 用于定义待办事项的优先级别
type Priority int

const (
	PriorityLow    Priority = 1 // 低优先级
	PriorityMedium Priority = 2 // 中优先级
	PriorityHigh   Priority = 3 // 高优先级
)

// Todo 待办事项模型
// 存储待办事项的详细信息，包括标题、描述、完成状态、优先级等
// 通过外键关联用户和分类信息
type Todo struct {
	Base
	Title       string     `json:"title" gorm:"size:128;not null;index"`                  // 待办事项标题，不超过128字符
	Description string     `json:"description" gorm:"size:1024"`                    // 待办事项描述，不超过1024字符
	Completed   bool       `json:"completed" gorm:"default:false;index"`                  // 完成状态，默认为未完成
	Priority    Priority   `json:"priority" gorm:"default:2;index"`                       // 优先级，默认为中优先级
	UserID      uint       `json:"userId" gorm:"not null;index"`             // 所属用户ID
	User        User       `gorm:"foreignKey:UserID" json:"-"`                      // 关联的用户信息，json序列化时忽略
	CategoryID  *uint      `json:"categoryId" gorm:"index"`                 // 所属分类ID，允许为空
	Category    *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"` // 关联的分类信息
	Reminders   []Reminder `json:"reminders,omitempty" gorm:"foreignKey:TodoID"`    // 关联的提醒列表
}
