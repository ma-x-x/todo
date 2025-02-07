package main

// 导入所需的包
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "todo/docs" // 导入swagger文档，用于API文档生成
	"todo/internal/models"
	routes "todo/internal/router"
	"todo/internal/service"
	"todo/pkg/cache"
	"todo/pkg/config"
	"todo/pkg/database"
	"todo/pkg/logger"
	"todo/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Swagger文档注解
// @title Todo API
// @version 1.0
// @description 这是一个待办事项管理系统的API服务
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @produce application/json
// @consume application/json

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 在Authorization头部输入"Bearer "后跟JWT令牌

// 程序入口函数
func main() {
	if err := run(); err != nil {
		log.Fatalf("***应用程序 启动失败***: %v", err)
	}
}

// run 应用程序的主运行逻辑
// 包含了完整的应用初始化和服务启动流程
func run() error {
	// 1. 加载配置文件
	// 从配置文件中读取应用所需的各项配置
	cfg, err := config.LoadConfig()
	// 输出配置
	fmt.Errorf("配置文件: %+v\n", cfg)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 2. 初始化日志系统
	// 设置日志记录器，用于记录应用运行时的各种信息
	if err := logger.Init(cfg.Logger); err != nil {
		return fmt.Errorf("初始化日志系统失败: %w", err)
	}

	// 3. 初始化数据库连接
	// 连接MySQL数据库，用于存储应用数据
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 验证数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 在初始化数据库连接后添加
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Category{}, &models.Reminder{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 验证索引是否存在
	for _, model := range []string{"users", "todos", "categories", "reminders"} {
		var count int64
		if err := db.Raw(`
			SELECT count(*) 
			FROM information_schema.statistics 
			WHERE table_schema = ? 
			AND table_name = ?`, 
			cfg.MySQL.Database, model).Count(&count).Error; err != nil {
			log.Printf("警告: 检查表 %s 的索引时出错: %v", model, err)
		} else if count == 0 {
			log.Printf("警告: 表 %s 可能缺少必要的索引", model)
		}
	}

	// 4. 初始化Redis缓存
	// 连接Redis，用于缓存和速率限制等功能
	rdb, err := cache.InitRedis(&cfg.Redis)
	if err != nil {
		return fmt.Errorf("初始化Redis失败: %w", err)
	}

	// 5. 初始化各个服务
	// 创建认证、待办事项、分类、提醒等服务的实例
	services := initServices(db, rdb, &cfg.JWT)

	// 6. 设置Gin框架的运行模式
	// 可以是debug或release模式
	if cfg.Server.Mode != "debug" && cfg.Server.Mode != "release" && cfg.Server.Mode != "test" {
		log.Printf("警告: 未知的服务器模式 '%s'，使用默认的 'release' 模式", cfg.Server.Mode)
		cfg.Server.Mode = "release"
	}
	gin.SetMode(cfg.Server.Mode)

	// 7. 初始化Web服务器
	r := gin.New()
	// 添加全局中间件
	r.Use(gin.Logger())                                  // 请求日志
	r.Use(gin.Recovery())                                // 错误恢复
	r.Use(middleware.RateLimiter(rdb, 100, time.Minute)) // Redis限流器，限制API访问频率

	// 初始化路由
	// 设置所有的API路由规则
	r = routes.InitRouter(cfg, services.auth, services.todo, services.category, services.reminder)

	// 8. 配置HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 在后台启动服务器
	go func() {
		log.Printf("服务器正在启动，地址：http://localhost:%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置5秒的超时时间来处理剩余请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	// 关闭Redis连接
	if err := cache.Close(); err != nil {
		log.Printf("关闭Redis连接失败: %v", err)
	}

	log.Println("服务器已成功关闭")
	return nil
}

// services 结构体用于组织所有服务实例
type services struct {
	auth     service.AuthService     // 认证服务
	todo     service.TodoService     // 待办事项服务
	category service.CategoryService // 分类服务
	reminder service.ReminderService // 提醒服务
}

// initServices 初始化所有服务
// 创建并返回各个服务的实例
func initServices(db *gorm.DB, rdb *redis.Client, jwtCfg *config.JWTConfig) *services {
	return &services{
		auth:     service.NewAuthService(db, rdb, jwtCfg),
		todo:     service.NewTodoService(db),
		category: service.NewCategoryService(db),
		reminder: service.NewReminderService(db),
	}
}
