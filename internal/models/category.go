package models

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
// @Description 待办事项分类
// 用于对待办事项进行分类管理
// 每个分类都属于特定用户，包含名称和颜色信息
type Category struct {
	Base
	// 分类名称
	Name string `json:"name" gorm:"size:100;not null;index"`
	// 分类描述
	Description string `json:"description" gorm:"size:500"`
	// 分类颜色
	Color string `json:"color" gorm:"size:7"`
	// 所属用户ID
	UserID uint `json:"userId" gorm:"not null;index"`
	// 关联的用户
	User User `json:"-" gorm:"foreignKey:UserID"`
	// 关联的待办事项列表
	Todos []Todo `json:"todos,omitempty" gorm:"foreignKey:CategoryID"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" example:"2024-02-08T17:08:41+08:00"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at" example:"2024-02-08T17:08:41+08:00"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Validate 验证分类数据
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("分类名称不能为空")
	}

	if len(c.Name) > 100 {
		return errors.New("分类名称不能超过100个字符")
	}

	if c.Description != "" && len(c.Description) > 500 {
		return errors.New("分类描述不能超过500个字符")
	}

	if c.Color != "" {
		// 验证颜色格式 #RRGGBB
		matched, _ := regexp.MatchString("^#[0-9A-Fa-f]{6}$", c.Color)
		if !matched {
			return errors.New("颜色格式无效，应为#RRGGBB格式")
		}
	}

	if c.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	return nil
}

// BeforeCreate 创建前的钩子函数
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if err := c.Validate(); err != nil {
		return err
	}
	return c.Base.BeforeCreate(tx)
}

// BeforeUpdate 更新前的钩子函数
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	if err := c.Validate(); err != nil {
		return err
	}
	return c.Base.BeforeUpdate(tx)
}
