// 新建 app 包来处理应用程序的核心逻辑
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo/internal/service"
	"todo/pkg/config"
	"todo/pkg/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App 结构体用于管理应用程序的核心组件
type App struct {
	cfg      *config.Config
	server   *http.Server
	services *service.ServiceCollection
	db       *gorm.DB      // 改为具体的 GORM DB 类型
	rdb      *redis.Client // 改为具体的 Redis 客户端类型
}

// New 创建新的应用程序实例
func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

// Initialize 初始化应用程序
func (a *App) Initialize() error {
	// 1. 初始化日志
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}
	log.Println("日志系统初始化成功")

	// 2. 初始化数据库
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	log.Println("数据库初始化成功")

	// 3. 初始化Redis
	if err := a.initRedis(); err != nil {
		return fmt.Errorf("初始化Redis失败: %w", err)
	}
	log.Println("Redis初始化成功")

	// 4. 初始化服务
	if err := a.initServices(); err != nil {
		return fmt.Errorf("初始化服务失败: %w", err)
	}
	log.Println("服务初始化成功")

	// 5. 初始化HTTP服务器
	if err := a.initHTTPServer(); err != nil {
		return fmt.Errorf("初始化HTTP服务器失败: %w", err)
	}
	log.Println("HTTP服务器初始化成功")

	return nil
}

// initLogger 初始化日志系统
func (a *App) initLogger() error {
	return logger.Init(a.cfg.Logger)
}

// Shutdown 优雅关闭应用程序
func (a *App) Shutdown(ctx context.Context) error {
	// 关闭 HTTP 服务器
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("关闭HTTP服务器失败: %w", err)
	}

	// TODO: 在这里添加其他资源的清理代码
	// 例如关闭数据库连接、Redis连接等

	return nil
}
