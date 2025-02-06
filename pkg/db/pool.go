package db

import (
	"fmt"
	"time"
	"todo/pkg/errors"
	"todo/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// DBConfig 数据库配置结构
// 支持主从分离的数据库连接配置
type DBConfig struct {
	Master      string        // 主库连接字符串，用于写操作
	Slaves      []string      // 从库连接字符串数组，用于读操作
	MaxIdle     int           // 最大空闲连接数
	MaxOpen     int           // 最大打开连接数
	MaxLifetime time.Duration // 连接最大生命周期
}

// NewDB 创建支持读写分离的数据库连接池
// 使用场景:
// 1. 大型应用或分布式系统
// 2. 需要读写分离的场景
// 3. 需要负载均衡的场景
// 4. 高并发访问场景
// 5. 对数据库访问性能要求较高的场景
func NewDB(config DBConfig) (*gorm.DB, error) {
	// 首先连接主库
	db, err := gorm.Open(mysql.Open(config.Master), &gorm.Config{})
	if err != nil {
		logger.Error().
			Str("dsn", config.Master).
			Err(err).
			Msg("主库连接失败")
		return nil, fmt.Errorf("%w: %v", errors.ErrDBConnection, err)
	}

	// 配置读写分离
	// Sources: 主库配置，用于写操作
	// Replicas: 从库配置，用于读操作
	// Policy: 从库负载均衡策略，这里使用随机策略
	resolverConfig := dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Master)},
		Replicas: make([]gorm.Dialector, len(config.Slaves)),
		Policy:   dbresolver.RandomPolicy{}, // 随机策略用于从库的负载均衡
	}

	// 配置所有从库连接
	for i, slave := range config.Slaves {
		resolverConfig.Replicas[i] = mysql.Open(slave)
		logger.Info().
			Int("index", i).
			Str("dsn", slave).
			Msg("配置从库连接")
	}

	// 注册 dbresolver 插件，启用读写分离
	err = db.Use(dbresolver.Register(resolverConfig))
	if err != nil {
		logger.Error().
			Err(err).
			Msg("注册读写分离插件失败")
		return nil, fmt.Errorf("配置读写分离失败: %w", err)
	}

	// 获取底层 *sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("获取底层数据库连接失败")
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 配置连接池
	// MaxIdleConns: 设置空闲连接池中的最大连接数
	// MaxOpenConns: 设置打开数据库连接的最大数量
	// ConnMaxLifetime: 设置连接可重用的最长时间
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	logger.Info().
		Str("master", config.Master).
		Int("slaves_count", len(config.Slaves)).
		Int("max_idle", config.MaxIdle).
		Int("max_open", config.MaxOpen).
		Msg("数据库连接池初始化成功")

	return db, nil
}
