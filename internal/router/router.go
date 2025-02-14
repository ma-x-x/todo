package router

import (
	"net/http"
	"time"
	"todo/api/v1/handlers"     // 处理 HTTP 请求的处理器
	"todo/internal/middleware" // 中间件组件
	"todo/internal/service"    // 业务逻辑服务
	"todo/pkg/config"          // 配置管理

	"github.com/gin-gonic/gin"     // Web 框架
	"github.com/redis/go-redis/v9" // Redis 客户端
)

// InitRouter 初始化所有路由配置
// 参数:
//   - cfg: 应用程序配置
//   - services: 所有服务的集合
//   - rdb: Redis 客户端实例
//   - r: Gin 引擎实例
//
// 路由结构:
//   - /api/v1/auth/: 认证相关路由（登录、注册）
//   - /api/v1/todos/: 待办事项管理
//   - /api/v1/categories/: 分类管理
//   - /api/v1/reminders/: 提醒管理
//   - /health: 健康检查接口
func InitRouter(cfg *config.Config, services *service.ServiceCollection, rdb *redis.Client, r *gin.Engine) {
	// 1. 注册全局中间件
	// 这些中间件会对所有请求生效
	r.Use(middleware.Cors())                // 处理跨域请求，必须第一个注册
	r.Use(gin.Recovery())                   // 从任何 panic 恢复，并返回 500 错误
	r.Use(middleware.RequestIDMiddleware()) // 为每个请求生成唯一 ID
	r.Use(middleware.LoggerMiddleware())    // 记录请求日志

	// 创建 API 版本路由组
	// 所有 API 路由都会以 /api/v1 开头
	api := r.Group("/api/v1")

	// 初始化所有 HTTP 处理器
	// 处理器负责处理具体的 HTTP 请求
	authHandler := handlers.NewAuthHandler(services.Auth)                    // 认证处理器
	todoHandler := handlers.NewTodoHandler(services.Todo, services.Category) // 待办事项处理器
	categoryHandler := handlers.NewCategoryHandler(services.Category)        // 分类处理器
	reminderHandler := handlers.NewReminderHandler(                          // 提醒处理器
		services.Reminder,
		services.Todo,
	)

	// 2. 注册不需要认证的路由
	// 这些路由可以直接访问，不需要登录
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register) // 用户注册
		auth.POST("/login", authHandler.Login)       // 用户登录
	}

	// 3. 注册需要认证的路由
	// 这些路由需要用户登录后才能访问
	authorized := api.Group("")
	authorized.Use(middleware.Auth(&cfg.JWT)) // 添加认证中间件
	{
		authorized.POST("/auth/logout", authHandler.Logout) // 用户登出

		// 待办事项相关路由
		todos := authorized.Group("/todos")
		{
			todos.POST("", todoHandler.Create)       // 创建待办事项
			todos.GET("", todoHandler.List)          // 获取待办事项列表
			todos.GET("/:id", todoHandler.Get)       // 获取单个待办事项
			todos.PUT("/:id", todoHandler.Update)    // 更新待办事项
			todos.DELETE("/:id", todoHandler.Delete) // 删除待办事项
		}

		// 分类相关路由
		categories := authorized.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)       // 创建分类
			categories.GET("", categoryHandler.List)          // 获取分类列表
			categories.GET("/:id", categoryHandler.Get)       // 获取单个分类
			categories.PUT("/:id", categoryHandler.Update)    // 更新分类
			categories.DELETE("/:id", categoryHandler.Delete) // 删除分类
		}

		// 提醒相关路由
		reminders := authorized.Group("/reminders")
		{
			reminders.POST("", reminderHandler.Create)            // 创建提醒
			reminders.GET("/todo/:todo_id", reminderHandler.List) // 获取某个待办事项的所有提醒
			reminders.GET("/:id", reminderHandler.Get)            // 获取单个提醒
			reminders.PUT("/:id", reminderHandler.Update)         // 更新提醒
			reminders.DELETE("/:id", reminderHandler.Delete)      // 删除提醒
		}
	}

	// 4. 健康检查路由
	// 用于监控系统状态，确保服务正常运行
	r.GET("/health", func(c *gin.Context) {
		// 检查数据库连接是否正常
		if err := services.CheckDatabaseHealth(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
			return
		}

		// 检查 Redis 连接是否正常
		if err := services.CheckRedisHealth(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "redis connection failed",
			})
			return
		}

		// 所有组件正常，返回健康状态
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339), // 返回当前时间，格式为 ISO8601
		})
	})
}
