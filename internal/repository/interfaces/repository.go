// Package interfaces 定义仓储层接口
package interfaces

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository 通用仓储接口
type Repository interface {
	// GetDB 获取数据库连接
	GetDB(ctx context.Context) *gorm.DB
}

// Repositories 所有仓储的集合
type Repositories struct {
	User     UserRepository
	Auth     AuthRepository
	Todo     TodoRepository
	Category CategoryRepository
	Reminder ReminderRepository
}

// BaseRepository 基础仓储接口
type BaseRepository interface {
	// Transaction 事务处理
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error

	// GetDB 获取数据库连接
	GetDB(ctx context.Context) *gorm.DB

	// Create 通用创建方法
	Create(ctx context.Context, model interface{}) error

	// Update 通用更新方法
	Update(ctx context.Context, model interface{}) error

	// Delete 通用删除方法
	Delete(ctx context.Context, model interface{}) error

	// GetByID 通用根据ID获取方法
	GetByID(ctx context.Context, id uint, model interface{}) error

	// List 通用列表查询方法
	List(ctx context.Context, offset, limit int, models interface{}, conditions ...interface{}) error

	// Count 通用计数方法
	Count(ctx context.Context, model interface{}, conditions ...interface{}) (int64, error)
}

// Options 仓储配置选项
type Options struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// NewRepositoriesFunc 定义创建仓储的函数类型
type NewRepositoriesFunc func(opts *Options) *Repositories
