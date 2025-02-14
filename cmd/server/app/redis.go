package app

import (
	"fmt"
	"todo/pkg/cache"
)

// initRedis 初始化Redis连接
func (a *App) initRedis() error {
	rdb, err := cache.InitRedis(&a.cfg.Redis)
	if err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	a.rdb = rdb
	return nil
}
