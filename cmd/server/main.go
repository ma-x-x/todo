package main

// 导入所需的包
import (
	// 标准库导入
	"context"   // 用于处理上下文，比如超时控制
	"fmt"       // 用于格式化输出
	"log"       // 用于日志记录
	"net/http"  // 提供 HTTP 客户端和服务端的实现
	"os"        // 提供操作系统功能
	"os/signal" // 用于处理操作系统信号
	"strings"   // 提供字符串操作函数
	"syscall"   // 提供系统调用功能
	"time"      // 处理时间相关操作

	// 项目内部包导入
	"todo/docs"                   // Swagger API 文档
	"todo/internal/models"        // 数据模型定义
	"todo/internal/repository"    // 数据访问层
	routes "todo/internal/router" // 路由定义
	"todo/internal/service"       // 业务逻辑层
	"todo/pkg/cache"              // 缓存相关
	"todo/pkg/config"             // 配置相关
	"todo/pkg/database"           // 数据库相关
	"todo/pkg/logger"             // 日志相关

	// 第三方包导入
	"github.com/gin-gonic/gin"                 // Web 框架
	swaggerFiles "github.com/swaggo/files"     // Swagger 文件处理
	ginSwagger "github.com/swaggo/gin-swagger" // Gin 的 Swagger 集成
)

// Swagger API 文档注解
// @title Todo API
// @version 1.0
// @description 这是一个待办事项管理系统的API服务
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @produce application/json
// @consume application/json

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请求头需要添加Bearer token

// main 函数是程序的入口点
// 它调用 run 函数启动服务器，如果发生错误则记录错误并退出程序
func main() {
	if err := run(); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}

// run 函数包含了应用程序的主要逻辑
// 它按照以下步骤初始化和启动服务：
// 1. 加载配置
// 2. 初始化日志系统
// 3. 连接数据库
// 4. 连接 Redis
// 5. 初始化各种服务
// 6. 设置并启动 Web 服务器
// 7. 处理优雅关闭
func run() error {
	// 1. 加载配置文件
	// 从配置文件（如 config.yaml）中读取应用程序的配置信息
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 2. 初始化日志系统
	// 设置日志记录器，用于记录应用运行时的各种信息
	if err := logger.Init(cfg.Logger); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	// 3. 初始化数据库连接
	// 连接 MySQL 数据库，这是应用程序的主要数据存储
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试数据库连接是否正常
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 使用 Ping 确保数据库连接正常
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 自动迁移数据库表结构
	// 这会自动创建或更新数据库表以匹配我们的模型定义
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Category{}, &models.Reminder{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 检查数据库表的索引是否存在
	// 索引能够提高数据库查询性能
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

	// 4. 初始化 Redis 连接
	// Redis 用于缓存和会话管理
	rdb, err := cache.InitRedis(&cfg.Redis)
	if err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}
	defer cache.Close() // 确保在程序结束时关闭 Redis 连接

	// 初始化各层的组件
	// 仓储层：负责数据访问
	repos := repository.NewRepositories(db, rdb)

	// 服务层：负责业务逻辑
	// 创建各种服务的实例
	authService := service.NewAuthService(repos.User, repos.Auth, &cfg.JWT)
	todoService := service.NewTodoService(repos.Todo)
	categoryService := service.NewCategoryService(repos.Category)
	reminderService := service.NewReminderService(repos.Reminder, repos.Todo)

	// 将所有服务集中管理
	services := service.NewServiceCollection(
		authService,
		todoService,
		categoryService,
		reminderService,
		db,  // 数据库连接
		rdb, // Redis 连接
	)

	// 5. 设置 Gin 框架的运行模式
	// debug: 开发模式，提供详细日志
	// release: 生产模式，性能优化
	// test: 测试模式
	log.Printf("设置 Gin 模式之前: %s", cfg.Server.Mode)
	if cfg.Server.Mode != "debug" && cfg.Server.Mode != "release" && cfg.Server.Mode != "test" {
		log.Printf("警告: 未知的服务器模式 '%s'，使用默认的 'release' 模式", cfg.Server.Mode)
		cfg.Server.Mode = "release"
	}
	log.Printf("最终使用的 Gin 模式: %s", cfg.Server.Mode)
	gin.SetMode(cfg.Server.Mode)

	// 6. 初始化 Web 服务器
	r := gin.New() // 创建一个新的 Gin 引擎实例

	// 设置路由
	routes.InitRouter(cfg, services, rdb, r)

	// 配置 Swagger API 文档
	// Swagger 提供了一个交互式的 API 文档界面
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置 Swagger 文档信息
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "Todo 应用后端 API 文档"
	docs.SwaggerInfo.Version = "1.0"

	// 设置 Swagger 的主机地址
	if cfg.Server.Mode == "release" && cfg.Server.SwaggerHost != "" {
		if !strings.Contains(cfg.Server.SwaggerHost, ":") {
			docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.SwaggerHost, cfg.Server.Port)
		} else {
			docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
		}
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Server.Port)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 7. 配置 HTTP 服务器
	// 设置超时时间和监听地址
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  60 * time.Second,  // 读取请求的超时时间
		WriteTimeout: 60 * time.Second,  // 写入响应的超时时间
		IdleTimeout:  120 * time.Second, // 空闲连接的超时时间
	}

	// 在后台启动服务器
	go func() {
		log.Printf("服务器启动在 http://localhost:%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v", err)
		}
	}()

	// 等待中断信号
	// 创建一个通道来接收操作系统的信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT（Ctrl+C）和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞直到收到信号
	log.Println("正在关闭服务器...")

	// 设置 5 秒的超时时间来处理剩余请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	// 这样可以确保正在处理的请求能够完成
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	log.Println("服务器已成功关闭")
	return nil
}
