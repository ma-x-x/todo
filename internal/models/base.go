package models

import (
	"time"

	"gorm.io/gorm"
)

// Base 基础模型
// 作为基础结构体，被其他模型继承，提供了通用的数据库字段
// 包含ID、创建时间、更新时间和软删除时间等基础字段
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`                    // 主键ID
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`     // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`     // 更新时间
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;index"` // 软删除时间
}
