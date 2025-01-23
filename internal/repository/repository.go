// Package repository 实现数据访问层，负责与数据库交互的所有操作
package repository

import (
	"gorm.io/gorm"
)

// NewUserRepository 创建用户仓储实例
// db: 数据库连接实例
// 返回: UserRepository 接口实现
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// NewTodoRepository 创建待办事项仓储实例
// db: 数据库连接实例
// 返回: TodoRepository 接口实现
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepo{db: db}
}

// NewCategoryRepository 创建分类仓储实例
// db: 数据库连接实例
// 返回: CategoryRepository 接口实现
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

// NewReminderRepository 创建提醒事项仓储实例
// db: 数据库连接实例
// 返回: ReminderRepository 接口实现
func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepo{db: db}
}
