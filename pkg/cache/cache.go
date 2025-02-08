package cache

import (
	"context"
	"fmt"
	"sync"
	"todo/pkg/config"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// InitRedis 初始化Redis客户端
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	var initErr error
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
			PoolSize: cfg.PoolSize,
		})

		// 测试连接
		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("连接Redis失败: %w", err)
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return rdb, nil
}

// GetRedis 获取全局Redis客户端实例
func GetRedis() *redis.Client {
	return rdb
}

// Close 关闭Redis连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
