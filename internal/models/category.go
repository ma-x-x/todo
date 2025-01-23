package models

// Category 分类模型
// @Description 待办事项分类
type Category struct {
	Base
	Name   string `json:"name" gorm:"size:32;not null"`
	Color  string `json:"color" gorm:"size:7"`
	UserID uint   `json:"user_id" gorm:"not null"`
}
