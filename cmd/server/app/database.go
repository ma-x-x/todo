package app

import (
	"fmt"
	"log"
	"todo/internal/models"
	"todo/pkg/database"

	"gorm.io/gorm"
)

// initDatabase 初始化数据库连接和迁移
func (a *App) initDatabase() error {
	// 连接数据库
	db, err := database.NewMySQLDB(a.cfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.Category{},
		&models.Reminder{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 检查数据库索引
	if err := a.checkDatabaseIndexes(db); err != nil {
		log.Printf("警告: 检查数据库索引时出现问题: %v", err)
	}

	a.db = db
	return nil
}

// checkDatabaseIndexes 检查数据库索引
func (a *App) checkDatabaseIndexes(db *gorm.DB) error {
	// 检查每个表的索引
	for _, model := range []string{"users", "todos", "categories", "reminders"} {
		var count int64
		if err := db.Raw(`
			SELECT count(*) 
			FROM information_schema.statistics 
			WHERE table_schema = ? 
			AND table_name = ?`,
			a.cfg.MySQL.Database, model).Count(&count).Error; err != nil {
			log.Printf("警告: 检查表 %s 的索引时出错: %v", model, err)
		} else if count == 0 {
			log.Printf("警告: 表 %s 可能缺少必要的索引", model)
		}
	}

	return nil
}
