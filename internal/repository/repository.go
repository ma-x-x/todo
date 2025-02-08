package repository

import (
	"todo/internal/repository/impl"
	"todo/internal/repository/interfaces"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewRepositories 创建仓储集合实例
func NewRepositories(db *gorm.DB, rdb *redis.Client) *interfaces.Repositories {
	opts := &interfaces.Options{
		DB:  db,
		RDB: rdb,
	}

	return &interfaces.Repositories{
		User:     impl.NewUserRepository(opts.DB),
		Auth:     impl.NewAuthRepository(opts.DB, opts.RDB),
		Todo:     impl.NewTodoRepository(opts.DB),
		Category: impl.NewCategoryRepository(opts.DB),
		Reminder: impl.NewReminderRepository(opts.DB),
	}
}

// WithTx 在事务中执行函数
func WithTx(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after Rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
