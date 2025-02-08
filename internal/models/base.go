package models

import (
	"time"

	"gorm.io/gorm"
)

// Base 基础模型
// 作为基础结构体，被其他模型继承，提供了通用的数据库字段
// 包含ID、创建时间、更新时间和软删除时间等基础字段
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`               // 主键ID
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"` // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`   // 软删除时间
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	// 按照依赖顺序进行迁移
	models := []interface{}{
		&User{},     // 用户表最先创建
		&Category{}, // 分类依赖用户表
		&Todo{},     // 待办事项依赖用户表和分类表
		&Reminder{}, // 提醒依赖待办事项表
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

// BeforeCreate 创建前钩子
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
