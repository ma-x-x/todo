package models

import (
	"time"

	"gorm.io/gorm"
)

// Base 模型定义，对应 gorm.Model
// @Description 基础模型
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
