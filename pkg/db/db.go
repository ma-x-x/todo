package db

import (
	"fmt"
	"todo/pkg/config"
	"todo/pkg/errors"
	"todo/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db 全局数据库连接实例
// 适用于单数据库连接的简单场景
var db *gorm.DB

// Init 初始化数据库连接
// 使用场景:
// 1. 小型应用或单体服务
// 2. 单一数据库连接
// 3. 不需要读写分离的场景
// 4. 对连接池有基本配置需求的场景
func Init(cfg *config.MySQLConfig) (*gorm.DB, error) {
	// 构建 MySQL DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// 建立数据库连接
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error().
			Err(err).
			Str("host", cfg.Host).
			Int("port", cfg.Port).
			Str("database", cfg.Database).
			Msg("数据库连接失败")
		return nil, fmt.Errorf("%w: %v", errors.ErrDBConnection, err)
	}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("获取底层数据库连接失败")
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 配置连接池参数
	// MaxIdleConns: 连接池中最大空闲连接数
	// MaxOpenConns: 连接池最大连接数
	// ConnMaxLifetime: 连接可复用的最大时间
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Int("maxIdleConns", cfg.MaxIdleConns).
		Int("maxOpenConns", cfg.MaxOpenConns).
		Msg("数据库连接成功")

	return db, nil
}

// GetDB 获取全局数据库连接实例
// 注意: 使用前必须先调用 Init 进行初始化
func GetDB() *gorm.DB {
	if db == nil {
		logger.Error().
			Msg("数据库未初始化，请先调用 Init 函数")
		panic("数据库未初始化")
	}
	return db
}
