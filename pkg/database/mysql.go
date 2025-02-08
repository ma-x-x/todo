package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"todo/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// DBConfig 数据库配置
type DBConfig struct {
	Master      string
	Slaves      []string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

// NewMySQLDB 创建MySQL连接(支持单库和主从)
func NewMySQLDB(cfg *config.Config) (*gorm.DB, error) {
	// 如果配置了从库，使用读写分离模式
	if len(cfg.MySQL.Slaves) > 0 {
		return NewDBWithReplicas(DBConfig{
			Master: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database),
			Slaves:      cfg.MySQL.Slaves,
			MaxIdle:     cfg.MySQL.MaxIdleConns,
			MaxOpen:     cfg.MySQL.MaxOpenConns,
			MaxLifetime: cfg.MySQL.ConnMaxLifetime,
		})
	}

	// 单库模式
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	)

	log.Printf("正在连接数据库: %s:%d/%s", cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(getLogLevel(cfg.Logger.Level)),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池
	configurePool(sqlDB, cfg.MySQL)
	return db, nil
}

// NewDBWithReplicas 创建支持读写分离的数据库连接
func NewDBWithReplicas(config DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Master), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 配置读写分离
	replicas := make([]gorm.Dialector, len(config.Slaves))
	for i, slave := range config.Slaves {
		replicas[i] = mysql.Open(slave)
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Master)},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

func configurePool(sqlDB *sql.DB, cfg config.MySQLConfig) {
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
}

// getLogLevel 根据配置的日志级别返回对应的gorm日志级别
func getLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "debug":
		return gormlogger.Info
	case "info":
		return gormlogger.Info
	case "warn":
		return gormlogger.Warn
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Silent
	}
}
