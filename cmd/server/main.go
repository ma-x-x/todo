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
	"todo/docs"
	_ "todo/docs" // 导入swagger文档，用于API文档生成
	"todo/internal/models"
	_ "todo/internal/models" // 导入模型定义
	"todo/internal/repository"
	routes "todo/internal/router"
	"todo/internal/service"
	"todo/pkg/cache"
	"todo/pkg/config"
	"todo/pkg/database"
	"todo/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger文档注解
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

// 程序入口函数
func main() {
	if err := run(); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}

// run 应用程序的主运行逻辑
// 包含了完整的应用初始化和服务启动流程
func run() error {
	// 1. 加载配置
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
	// 连接MySQL数据库，用于存储应用数据
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 验证数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 在初始化数据库连接后添加
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Category{}, &models.Reminder{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
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

	// 4. 初始化Redis连接
	rdb, err := cache.InitRedis(&cfg.Redis)
	if err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}
	defer cache.Close()

	// 初始化仓储层
	repos := repository.NewRepositories(db, rdb)

	// 初始化服务层
	authService := service.NewAuthService(repos.User, repos.Auth, &cfg.JWT)
	todoService := service.NewTodoService(repos.Todo)
	categoryService := service.NewCategoryService(repos.Category)
	reminderService := service.NewReminderService(repos.Reminder, repos.Todo)

	services := service.NewServiceCollection(
		authService,
		todoService,
		categoryService,
		reminderService,
		db,  // 数据库连接
		rdb, // Redis连接
	)

	// 5. 设置Gin框架的运行模式
	log.Printf("设置 Gin 模式之前: %s", cfg.Server.Mode)
	if cfg.Server.Mode != "debug" && cfg.Server.Mode != "release" && cfg.Server.Mode != "test" {
		log.Printf("警告: 未知的服务器模式 '%s'，使用默认的 'release' 模式", cfg.Server.Mode)
		cfg.Server.Mode = "release"
	}
	log.Printf("最终使用的 Gin 模式: %s", cfg.Server.Mode)
	gin.SetMode(cfg.Server.Mode)

	// 6. 初始化Web服务器
	r := gin.New()

	// 初始化路由
	routes.InitRouter(cfg, services, rdb, r)

	// 添加 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 配置 Swagger
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "Todo 应用后端 API 文档"
	docs.SwaggerInfo.Version = "1.0"

	// 使用配置中的 swagger_host
	if cfg.Server.Mode == "release" && cfg.Server.SwaggerHost != "" {
		docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Server.Port)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 7. 配置HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 在后台启动服务器
	go func() {
		log.Printf("服务器启动在 http://localhost:%d", cfg.Server.Port)
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

	log.Println("服务器已成功关闭")
	return nil
}
