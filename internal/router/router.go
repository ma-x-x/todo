package router

import (
	"net/http"
	"time"
	"todo/api/v1/handlers"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// InitRouter 初始化路由配置
func InitRouter(cfg *config.Config, services *service.ServiceCollection, rdb *redis.Client, r *gin.Engine) {
	// 1. 首先注册全局中间件
	r.Use(middleware.Cors()) // CORS 中间件必须第一个注册
	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// API 路由组
	api := r.Group("/api/v1")

	// 创建所有处理器实例
	authHandler := handlers.NewAuthHandler(services.Auth)
	todoHandler := handlers.NewTodoHandler(services.Todo, services.Category)
	categoryHandler := handlers.NewCategoryHandler(services.Category)
	reminderHandler := handlers.NewReminderHandler(
		services.Reminder,
		services.Todo,
	)

	// 2. 注册认证相关路由（不需要认证的路由）
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// 3. 需要认证的路由
	authorized := api.Group("")
	authorized.Use(middleware.Auth(&cfg.JWT))
	{
		authorized.POST("/auth/logout", authHandler.Logout)
		// 待办事项路由
		todos := authorized.Group("/todos")
		{
			todos.POST("", todoHandler.Create)
			todos.GET("", todoHandler.List)
			todos.GET("/:id", todoHandler.Get)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}

		// 分类路由
		categories := authorized.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.List)
			categories.GET("/:id", categoryHandler.Get)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// 提醒路由
		reminders := authorized.Group("/reminders")
		{
			reminders.POST("", reminderHandler.Create)
			reminders.GET("/todo/:todo_id", reminderHandler.List)
			reminders.GET("/:id", reminderHandler.Get)
			reminders.PUT("/:id", reminderHandler.Update)
			reminders.DELETE("/:id", reminderHandler.Delete)
		}
	}

	// 添加健康检查路由
	r.GET("/health", func(c *gin.Context) {
		// 检查数据库连接
		if err := services.CheckDatabaseHealth(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
			return
		}

		// 检查 Redis 连接
		if err := services.CheckRedisHealth(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "redis connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}
