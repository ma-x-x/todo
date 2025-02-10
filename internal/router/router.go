package router

import (
	"todo/api/v1/handlers"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// InitRouter 初始化路由配置
func InitRouter(cfg *config.Config, services *service.Services, rdb *redis.Client, r *gin.Engine) {
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
}
