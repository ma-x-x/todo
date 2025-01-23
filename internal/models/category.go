package models

// Category 分类模型
// 用于对待办事项进行分类管理
// 每个分类都属于特定用户，包含名称和颜色信息
type Category struct {
	Base
	Name   string `json:"name" gorm:"size:32;not null"` // 分类名称，不超过32字符
	Color  string `json:"color" gorm:"size:7"`          // 分类颜色，使用十六进制颜色码(如 #FF0000)
	UserID uint   `json:"user_id" gorm:"not null"`      // 所属用户ID
}
