package routes

import (
	"todo-demo/internal/middleware"
	"todo-demo/internal/service"
	"todo-demo/pkg/config"

	"todo-demo/api/v1/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
// 该函数负责设置所有的HTTP路由规则，包括API端点、中间件和Swagger文档
func InitRouter(cfg *config.Config, authService service.AuthService, todoService service.TodoService,
	categoryService service.CategoryService, reminderService service.ReminderService) *gin.Engine {

	// 创建一个新的Gin引擎实例
	r := gin.New()

	// 配置全局中间件
	r.Use(gin.Logger())                  // 请求日志记录
	r.Use(gin.Recovery())                // 错误恢复，防止服务器崩溃
	r.Use(middleware.LoggerMiddleware()) // 自定义日志中间件
	r.Use(middleware.CORSMiddleware())   // 跨域资源共享(CORS)支持

	// Swagger API文档路由
	// 访问 /swagger/index.html 可以查看API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 版本分组
	v1 := r.Group("/api/v1")
	{
		// 健康检查端点
		// 用于监控服务是否正常运行
		v1.GET("/health", handlers.Health)

		// 认证相关路由组
		// 包含注册和登录功能
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register(authService)) // 用户注册
			auth.POST("/login", handlers.Login(authService))       // 用户登录
		}

		// 需要认证的路由组
		// 以下所有路由都需要有效的JWT令牌才能访问
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware(&cfg.JWT))
		{
			// 待办事项管理路由组
			todos := authorized.Group("/todos")
			{
				todos.POST("", handlers.CreateTodo(todoService))       // 创建待办事项
				todos.GET("", handlers.ListTodos(todoService))         // 获取待办事项列表
				todos.GET("/:id", handlers.GetTodo(todoService))       // 获取单个待办事项
				todos.PUT("/:id", handlers.UpdateTodo(todoService))    // 更新待办事项
				todos.DELETE("/:id", handlers.DeleteTodo(todoService)) // 删除待办事项
			}

			// 分类管理路由组
			categories := authorized.Group("/categories")
			{
				categories.POST("", handlers.CreateCategory(categoryService))       // 创建分类
				categories.GET("", handlers.ListCategories(categoryService))        // 获取分类列表
				categories.PUT("/:id", handlers.UpdateCategory(categoryService))    // 更新分类
				categories.DELETE("/:id", handlers.DeleteCategory(categoryService)) // 删除分类
			}

			// 提醒管理路由组
			reminders := authorized.Group("/reminders")
			{
				reminders.POST("", handlers.CreateReminder(reminderService))             // 创建提醒
				reminders.GET("/todo/:todo_id", handlers.ListReminders(reminderService)) // 获取待办事项的提醒列表
				reminders.PUT("/:id", handlers.UpdateReminder(reminderService))          // 更新提醒
				reminders.DELETE("/:id", handlers.DeleteReminder(reminderService))       // 删除提醒
			}
		}
	}

	return r
}
