# Todo API 项目开发指南

## 目录

1. [项目学习指南](#1-项目学习指南)
2. [开发环境准备](#2-开发环境准备) 
3. [项目目录规范](#3-项目目录规范)
4. [Go 编码规范](#4-go-编码规范)
5. [项目开发流程详解](#5-项目开发流程详解)
6. [实践指南](#6-实践指南)
7. [调试与优化](#7-调试与优化)
8. [常见问题解决](#8-常见问题解决)
9. [API 开发与测试](#9-api-开发与测试)
10. [高级特性实现指南](#10-高级特性实现指南)
11. [项目最佳实践总结](#11-项目最佳实践总结)

## 1. 项目学习指南

### 1.1 学习路径

#### 1.1.1 基础知识准备
1. Go 语言基础
   - 语法基础：变量、控制流、函数
   - 并发编程：goroutine、channel
   - 接口和面向对象
   - 推荐资源：[Go by Example](https://gobyexample.com/)

2. Web 开发基础
   - HTTP 协议
   - RESTful API 设计
   - JSON 数据格式
   - 数据库基础

#### 1.1.2 核心概念学习
1. 项目架构
   ```
   表现层 (API) -> 业务层 (Service) -> 数据访问层 (Repository) -> 数据库
   ```
   - 每一层的职责
   - 层与层之间如何交互
   - 为什么要分层

2. 依赖注入
   ```go
   // 不使用依赖注入
   type UserService struct {
       repo *UserRepository  // 直接依赖具体实现
   }

   // 使用依赖注入
   type UserService struct {
       repo UserRepository   // 依赖接口
   }
   ```

3. 中间件机制
   ```go
   func LoggerMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           start := time.Now()
           c.Next()
           duration := time.Since(start)
           log.Printf("请求处理时间: %v", duration)
       }
   }
   ```

### 1.3 核心概念详解

#### 1.3.1 分层架构详解
```
┌─────────────┐
│    API层    │ 处理HTTP请求，参数验证，返回响应
├─────────────┤
│   Service层 │ 实现业务逻辑，事务管理，数据组装
├─────────────┤
│ Repository层│ 数据访问，CRUD操作，查询优化
├─────────────┤
│   Model层   │ 数据模型定义，字段验证，关联关系
└─────────────┘
```

1. API层（表现层）
   - 职责：接收请求，返回响应
   - 主要工作：
     * 请求参数解析和验证
     * 调用相应的服务层方法
     * 处理响应结果
     * 错误处理和响应封装
   - 代码示例：
   ```go
   func (h *TodoHandler) Create(c *gin.Context) {
       // 1. 参数验证
       var req CreateTodoRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           c.JSON(400, ErrorResponse{Error: "无效的请求参数"})
           return
       }

       // 2. 调用服务层
       todo, err := h.todoService.Create(req)
       if err != nil {
           c.JSON(500, ErrorResponse{Error: "创建失败"})
           return
       }

       // 3. 返回结果
       c.JSON(200, todo)
   }
   ```

2. Service层（业务层）
   - 职责：实现业务逻辑
   - 主要工作：
     * 业务规则验证
     * 数据处理和转换
     * 事务管理
     * 调用其他服务
   - 代码示例：
   ```go
   func (s *TodoService) Create(req CreateTodoRequest) (*Todo, error) {
       // 1. 业务规则验证
       if err := s.validateTodo(req); err != nil {
           return nil, err
       }

       // 2. 开启事务
       tx := s.db.Begin()
       defer func() {
           if r := recover(); r != nil {
               tx.Rollback()
           }
       }()

       // 3. 创建待办事项
       todo := &Todo{
           Title:       req.Title,
           Description: req.Description,
           UserID:      req.UserID,
       }

       if err := tx.Create(todo).Error; err != nil {
           tx.Rollback()
           return nil, err
       }

       // 4. 提交事务
       if err := tx.Commit().Error; err != nil {
           return nil, err
       }

       return todo, nil
   }
   ```

3. Repository层（数据访问层）
   - 职责：数据库操作封装
   - 主要工作：
     * CRUD 操作实现
     * 查询优化
     * 数据映射
     * 缓存处理
   - 代码示例：
   ```go
   func (r *TodoRepository) Create(todo *Todo) error {
       // 1. 检查缓存
       if cached, err := r.cache.Get(todo.ID); err == nil {
           return nil
       }

       // 2. 写入数据库
       if err := r.db.Create(todo).Error; err != nil {
           return err
       }

       // 3. 更新缓存
       r.cache.Set(todo.ID, todo, time.Hour)
       return nil
   }
   ```

4. Model层（数据模型层）
   - 职责：定义数据结构
   - 主要工作：
     * 字段定义和验证
     * 关联关系定义
     * 模型方法实现
     * 钩子函数定义
   - 代码示例：
   ```go
   type Todo struct {
       gorm.Model
       Title       string    `gorm:"size:128;not null" json:"title"`
       Description string    `gorm:"size:1024" json:"description"`
       DueDate     time.Time `gorm:"index" json:"due_date"`
       Status      string    `gorm:"size:20;default:'pending'" json:"status"`
       UserID      uint      `gorm:"index;not null" json:"user_id"`
       User        User      `gorm:"foreignKey:UserID" json:"user"`
   }

   // 验证方法
   func (t *Todo) Validate() error {
       if strings.TrimSpace(t.Title) == "" {
           return errors.New("标题不能为空")
       }
       if t.DueDate.Before(time.Now()) {
           return errors.New("截止日期不能早于当前时间")
       }
       return nil
   }

   // 钩子函数
   func (t *Todo) BeforeCreate(tx *gorm.DB) error {
       return t.Validate()
   }
   ```

#### 1.3.2 依赖注入详解

1. 什么是依赖注入？
   - 定义：一种设计模式，通过外部注入依赖而不是在内部创建
   - 目的：降低耦合度，提高代码可测试性和可维护性
   - 优势：
     * 方便进行单元测试
     * 便于切换实现
     * 提高代码重用性

2. 依赖注入的方式：
   ```go
   // 1. 构造函数注入
   type TodoService struct {
       repo   TodoRepository
       cache  Cache
       logger Logger
   }

   func NewTodoService(repo TodoRepository, cache Cache, logger Logger) *TodoService {
       return &TodoService{
           repo:   repo,
           cache:  cache,
           logger: logger,
       }
   }

   // 2. 属性注入
   type TodoService struct {
       Repo   TodoRepository
       Cache  Cache
       Logger Logger
   }

   // 3. 接口注入
   type TodoServiceInjector interface {
       InjectTodoRepository(TodoRepository)
       InjectCache(Cache)
       InjectLogger(Logger)
   }
   ```

3. 依赖注入容器：
   ```go
   // 使用 wire 进行依赖注入
   // wire.go
   func InitializeTodoAPI() *TodoHandler {
       wire.Build(
           NewTodoHandler,
           NewTodoService,
           NewTodoRepository,
           NewCache,
           NewLogger,
       )
       return nil
   }
   ```

#### 1.3.3 中间件机制详解

1. 中间件的作用：
   - 横切关注点处理
   - 请求前预处理
   - 响应后处理
   - 错误处理

2. 中间件类型：
   ```go
   // 1. 认证中间件
   func AuthMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // 前置处理
           token := c.GetHeader("Authorization")
           if token == "" {
               c.AbortWithStatus(401)
               return
           }

           // 验证 token
           claims, err := validateToken(token)
           if err != nil {
               c.AbortWithStatus(401)
               return
           }

           // 设置用户信息
           c.Set("userID", claims.UserID)

           // 继续处理请求
           c.Next()

           // 后置处理
           // 可以在这里添加一些清理工作
       }
   }

   // 2. 日志中间件
   func LoggerMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // 开始时间
           start := time.Now()

           // 请求信息记录
           path := c.Request.URL.Path
           raw := c.Request.URL.RawQuery

           // 处理请求
           c.Next()

           // 结束时间
           end := time.Now()
           latency := end.Sub(start)

           // 记录日志
           log.Printf("[%s] %s?%s %v %d",
           c.Request.Method,
           path,
           raw,
           latency,
           c.Writer.Status(),
       )
       }
   }

   // 3. 错误处理中间件
   func ErrorMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           c.Next()

           // 检查是否有错误
           if len(c.Errors) > 0 {
               // 获取最后一个错误
               err := c.Errors.Last()

               // 错误类型判断
               switch e := err.Err.(type) {
               case *CustomError:
                   c.JSON(e.Status, gin.H{
                       "error": e.Message,
                   })
               default:
                   c.JSON(500, gin.H{
                       "error": "Internal Server Error",
                   })
               }
           }
       }
   }
   ```

3. 中间件链执行顺序：
   ```
   请求 → 中间件1前置 → 中间件2前置 → 处理函数 → 中间件2后置 → 中间件1后置 → 响应
   ```

## 2. 开发环境准备

### 2.1 必需环境
1. Go 1.21+
2. MySQL 8.0+
3. Git
4. Make

### 2.2 可选环境
1. Docker 20.10+
2. Redis 6.0+
3. Kubernetes 1.20+
4. IDE (推荐 GoLand 或 VSCode)

### 2.3 推荐工具
1. go-swagger: API 文档生成
2. golangci-lint: 代码检查
3. air: 热重载
4. mockgen: 测试 mock 生成

## 3. 项目目录规范

### 目录划分标准

1. `/api`：对外暴露的 API 接口定义
   - 按版本划分子目录（如 v1、v2）
   - 每个版本目录下包含：
     - `dto/`: 数据传输对象
     - `handlers/`: 请求处理器
     - `routes/`: 路由定义
   - 避免在此层引入业务逻辑

2. `/internal`：内部代码包
   - 不对外暴露的代码
   - 包含核心业务逻辑
   - 子目录职责必须单一
   - 遵循依赖注入原则

3. `/pkg`：可被外部应用程序使用的代码库
   - 通用工具和组件
   - 保持高内聚、低耦合
   - 完善的单元测试
   - 详细的使用文档

4. `/cmd`：项目主要的应用程序
   - 代码尽量简单
   - 仅包含初始化和启动逻辑
   - 通过依赖注入组装各个组件

5. `/configs`：配置文件目录
   - 区分环境的配置文件
   - 敏感信息使用环境变量
   - 配置项要有注释说明

6. `/test`：额外的测试应用程序和测试数据
   - 集成测试
   - 性能测试
   - 测试工具
   - 测试数据

### 代码分层规范

1. 表现层（API 层）
   ```
   /api
   ├── v1
   │   ├── dto          # 数据传输对象
   │   ├── handlers     # 请求处理
   │   └── routes       # 路由定义
   ```

2. 业务层（Service 层）
   ```
   /internal/service
   ├── interfaces     # 接口定义
   └── impl          # 接口实现
   ```

3. 数据访问层（Repository 层）
   ```
   /internal/repository
   ├── interfaces    # 仓储接口
   └── impl         # 接口实现
   ```

4. 领域模型层（Model 层）
   ```
   /internal/models
   ├── entities     # 实体定义
   └── valueobjects # 值对象
   ```

## 4. Go 编码规范

### 1. 命名规范

1. 包名
   - 使用小写单词
   - 简短且有意义
   - 避免下划线
   ```go
   // Good
   package models
   package repository
   
   // Bad
   package Models
   package user_repository
   ```

2. 文件名
   - 小写单词
   - 使用下划线分隔
   ```
   user_repository.go
   auth_middleware.go
   ```

3. 接口名
   - 通常以 er 结尾
   - 单一职责
   ```go
   type Reader interface { ... }
   type Writer interface { ... }
   type TodoRepository interface { ... }
   ```

4. 变量名
   - 驼峰命名
   - 简短清晰
   - 避免单字母（除循环计数器外）
   ```go
   // Good
   userID := 1
   todoList := []Todo{}
   
   // Bad
   UserId := 1
   tdl := []Todo{}
   ```

### 2. 代码组织

1. 文件结构
   ```go
   // 包声明
   package user
   
   // 导入
   import (
       "errors"
       "time"
       
       "github.com/gin-gonic/gin"
   )
   
   // 常量
   const maxRetries = 3
   
   // 类型定义
   type User struct { ... }
   
   // 接口定义
   type UserRepository interface { ... }
   
   // 函数/方法实现
   func (u *User) ValidatePassword() bool { ... }
   ```

2. 接口实现
   ```go
   // internal/repository/user.go
   type UserRepository interface {
       Create(user *models.User) error
       FindByID(id uint) (*models.User, error)
   }
   
   // internal/repository/impl/user.go
   type userRepository struct {
       db *gorm.DB
   }
   
   func NewUserRepository(db *gorm.DB) UserRepository {
       return &userRepository{db: db}
   }
   ```

### 3. 错误处理

1. 错误定义
   ```go
   // pkg/errors/errors.go
   var (
       ErrUserNotFound = errors.New("user not found")
       ErrInvalidInput = errors.New("invalid input")
   )
   ```

2. 错误处理
   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to create user: %w", err)
   }
   
   // Bad
   if err != nil {
       return err
   }
   ```

### 4. 注释规范

1. 包注释
   ```go
   // Package models 定义了应用程序的数据模型
   // 包括用户、待办事项等核心实体
   package models
   ```

2. 函数注释
   ```go
   // CreateUser 创建新用户
   // 如果用户名已存在，返回 ErrUserExists
   // 如果输入无效，返回 ErrInvalidInput
   func CreateUser(user *User) error {
       // ...
   }
   ```

3. 类型注释
   ```go
   // User 表示系统中的用户实体
   // 包含用户的基本信息和认证信息
   type User struct {
       ID       uint   `json:"id"`
       Username string `json:"username"`
   }
   ```

### 5. 测试规范

1. 测试文件命名
   ```
   user_test.go
   todo_service_test.go
   ```

2. 测试函数
   ```go
   func TestUserCreate(t *testing.T) {
       // 准备测试数据
       user := &models.User{
           Username: "testuser",
           Email:    "test@example.com",
       }
       
       // 执行测试
       err := userService.Create(user)
       
       // 断言结果
       assert.NoError(t, err)
       assert.NotZero(t, user.ID)
   }
   ```

3. 表驱动测试
   ```go
   func TestValidatePassword(t *testing.T) {
       tests := []struct {
           name     string
           password string
           want     bool
       }{
           {"valid password", "Password123!", true},
           {"too short", "Pwd1!", false},
           {"no number", "Password!", false},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got := ValidatePassword(tt.password)
               assert.Equal(t, tt.want, got)
           })
       }
   }
   ```

### 6. 并发处理

1. goroutine 使用
   ```go
   // Good
   go func() {
       defer func() {
           if err := recover(); err != nil {
               log.Printf("panic recovered: %v", err)
           }
       }()
       // 处理逻辑
   }()
   ```

2. channel 使用
   ```go
   // 带缓冲的 channel
   ch := make(chan Task, 100)
   
   // 使用 select 处理多个 channel
   select {
   case task := <-taskChan:
       processTask(task)
   case <-time.After(5 * time.Second):
       return ErrTimeout
   }
   ```

### 7. 性能优化

1. 内存分配
   ```go
   // Good - 预分配内存
   users := make([]User, 0, expectedSize)
   
   // Bad - 频繁扩容
   var users []User
   ```

2. 字符串拼接
   ```go
   // Good
   var builder strings.Builder
   builder.WriteString("Hello")
   builder.WriteString(" World")
   
   // Bad
   str := "Hello" + " " + "World"
   ```

### 8. 依赖注入

```go
// Good
type TodoService struct {
    repo   TodoRepository
    cache  Cache
    logger Logger
}

func NewTodoService(repo TodoRepository, cache Cache, logger Logger) *TodoService {
    return &TodoService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

// Bad
type TodoService struct {
    repo   *TodoRepository
    cache  *Redis
    logger *Logger
}
```

## 5. 项目开发流程详解

### 1. 项目初始化

#### 1.1 创建项目目录
```bash
# 创建项目根目录
mkdir todo
cd todo

# 创建项目结构
mkdir -p cmd/server
mkdir -p internal/{models,repository,service,middleware}
mkdir -p api/v1/{handlers,dto,routes}
mkdir -p pkg/{database,errors,logger,utils}
mkdir -p configs
mkdir -p deployments/{docker,kubernetes}
mkdir -p test/testutils
```

#### 1.2 初始化 Go 模块
```bash
go mod init todo
```

#### 1.3 创建 .gitignore
```bash
touch .gitignore
# 添加需要忽略的文件和目录
```

### 2. 配置管理

#### 2.1 创建配置文件
```yaml
# configs/config.yaml
server:
  mode: development
  port: 8080

database:
  host: localhost
  port: 3306
  username: root
  password: root
  dbname: todo_db

jwt:
  secret: your_jwt_secret_key
  expire_hours: 24
```

#### 2.2 实现配置加载
```go
// pkg/config/config.go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

// ... 配置结构体定义 ...

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 3. 数据库设计

#### 3.1 用户模型
```go
// internal/models/user.go
package models

import (
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;size:32"`
    Password string `gorm:"size:128"`
    Email    string `gorm:"size:128"`
}

func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
```

#### 3.2 待办事项模型
```go
// internal/models/todo.go
package models

import "gorm.io/gorm"

type Todo struct {
    gorm.Model
    Title       string `gorm:"size:128"`
    Description string `gorm:"size:1024"`
    Completed   bool   `gorm:"default:false"`
    UserID      uint
    User        User `gorm:"foreignKey:UserID"`
}
```

### 4. 数据库连接

```go
// pkg/database/mysql.go
package database

import (
    "fmt"
    "todo/internal/models"
    "todo/pkg/config"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
    )

    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    // 自动迁移
    return db.AutoMigrate(&models.User{}, &models.Todo{})
}

func GetDB() *gorm.DB {
    return db
}
```

### 5. 实现各层功能

#### 5.1 Repository 层
Repository 层负责数据访问，实现与数据库的交互。

```go
// internal/repository/user.go
package repository

import (
    "todo/internal/models"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
}

// ... 实现代码 ...

// internal/repository/todo.go
type TodoRepository interface {
    Create(todo *models.Todo) error
    Update(todo *models.Todo) error
    Delete(id uint) error
    FindByID(id uint) (*models.Todo, error)
    FindByUserID(userID uint) ([]models.Todo, error)
}

// ... 实现代码 ...
```

#### 5.2 Service 层
Service 层实现业务逻辑，处理数据验证和业务规则。

```go
// internal/service/auth.go
package service

type AuthService interface {
    Register(username, password, email string) error
    Login(username, password string) (string, error)
}

// internal/service/todo.go
type TodoService interface {
    Create(userID uint, title, description string) (*models.Todo, error)
    Update(id, userID uint, title, description string, completed *bool) (*models.Todo, error)
    Delete(id, userID uint) error
    GetByID(id, userID uint) (*models.Todo, error)
    ListByUserID(userID uint) ([]models.Todo, error)
}
```

#### 5.3 Handler 层
Handler 层处理 HTTP 请求，实现 API 接口。

```go
// api/v1/handlers/auth.go
package handlers

type AuthHandler struct {
    authService service.AuthService
}

func (h *AuthHandler) Register(c *gin.Context) {
    // 实现注册逻辑
}

func (h *AuthHandler) Login(c *gin.Context) {
    // 实现登录逻辑
}

// api/v1/handlers/todo.go
type TodoHandler struct {
    todoService service.TodoService
}

// ... 实现 CRUD 处理方法 ...
```

### 6. 中间件实现

#### 6.1 JWT 认证中间件
```go
// internal/middleware/auth.go
package middleware

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT 认证逻辑
    }
}
```

#### 6.2 日志中间件
```go
// internal/middleware/logger.go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 日志记录逻辑
    }
}
```

### 7. 路由配置

```go
// api/v1/routes/routes.go
package routes

func InitRouter() *gin.Engine {
    r := gin.New()
    
    // 使用中间件
    r.Use(middleware.LoggerMiddleware())
    r.Use(middleware.RecoveryMiddleware())
    
    // API 路由组
    v1 := r.Group("/api/v1")
    {
        // 认证路由
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
        }
        
        // 需要认证的路由
        todos := v1.Group("/todos")
        todos.Use(middleware.AuthMiddleware())
        {
            todos.POST("", todoHandler.Create)
            todos.GET("", todoHandler.List)
            todos.GET("/:id", todoHandler.Get)
            todos.PUT("/:id", todoHandler.Update)
            todos.DELETE("/:id", todoHandler.Delete)
        }
    }
    
    return r
}
```

### 8. 主程序入口

```go
// cmd/server/main.go
package main

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // 初始化数据库
    if err := database.Init(&cfg.Database); err != nil {
        log.Fatal(err)
    }
    
    // 初始化路由
    r := routes.InitRouter()
    
    // 启动服务器
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

### 9. 单元测试

编写各层的单元测试，确保代码质量：

```go
// internal/repository/user_test.go
// internal/service/auth_test.go
// api/v1/handlers/todo_test.go
```

### 10. API 文档

使用 Swagger 生成 API 文档：

```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go
```

### 11. 容器化

#### 11.1 Dockerfile
```dockerfile
# deployments/docker/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/bin/server .
EXPOSE 8080
CMD ["./server"]
```

#### 11.2 Docker Compose
```yaml
# deployments/docker/docker-compose.yml
version: "3.8"
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mysql
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: todo_db
```

### 12. Kubernetes 部署

准备 Kubernetes 配置文件：
- configmap.yaml
- secret.yaml
- deployment.yaml
- service.yaml
- ingress.yaml

## 6. 实践指南

### 1. 循序渐进的开发步骤

1. 基础框架搭建
   ```go
   // 1. 创建主程序入口
   func main() {
       // 初始化配置
       config.Init()
       
       // 初始化数据库
       db.Init()
       
       // 创建路由
       r := gin.Default()
       
       // 注册路由
       routes.Register(r)
       
       // 启动服务器
       r.Run()
   }
   ```

2. 实现用户认证
   ```go
   // 2. 实现注册功能
   func (h *AuthHandler) Register(c *gin.Context) {
       // 1) 解析请求数据
       var req RegisterRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           // 处理错误
           return
       }
       
       // 2) 业务逻辑处理
       err := h.authService.Register(req)
       if err != nil {
           // 处理错误
           return
       }
       
       // 3) 返回结果
       c.JSON(http.StatusOK, gin.H{
           "message": "注册成功",
       })
   }
   ```

3. 实现待办事项功能
   ```go
   // 3. 实现创建待办事项
   func (h *TodoHandler) Create(c *gin.Context) {
       // 1) 获取当前用户
       userID := auth.GetUserID(c)
       
       // 2) 解析请求数据
       var req CreateTodoRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           // 处理错误
           return
       }
       
       // 3) 创建待办事项
       todo, err := h.todoService.Create(userID, req)
       if err != nil {
           // 处理错误
           return
       }
       
       // 4) 返回结果
       c.JSON(http.StatusOK, todo)
   }
   ```

#### 2.2 关键功能实现解析

1. JWT 认证实现
   ```go
   // JWT 认证中间件
   func AuthMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // 1. 获取 token
           token := c.GetHeader("Authorization")
           
           // 2. 验证 token
           claims, err := jwt.ValidateToken(token)
           if err != nil {
               c.AbortWithStatus(http.StatusUnauthorized)
               return
           }
           
           // 3. 设置用户信息到上下文
           c.Set("userID", claims.UserID)
           
           c.Next()
       }
   }
   ```

2. 数据验证实现
   ```go
   // 请求数据验证
   type CreateTodoRequest struct {
       Title       string `json:"title" binding:"required,min=1,max=128"`
       Description string `json:"description" binding:"max=1024"`
   }
   
   // 自定义验证器
   func validateTodoRequest(req CreateTodoRequest) error {
       if strings.TrimSpace(req.Title) == "" {
           return errors.New("标题不能为空")
       }
       return nil
   }
   ```

3. 错误处理实现
   ```go
   // 统一错误处理
   func ErrorHandler() gin.HandlerFunc {
       return func(c *gin.Context) {
           c.Next()
           
           // 检查是否有错误
           if len(c.Errors) > 0 {
               err := c.Errors.Last()
               
               // 转换为业务错误
               var bizErr *errors.BusinessError
               if errors.As(err, &bizErr) {
                   c.JSON(bizErr.Status(), gin.H{
                       "error": bizErr.Error(),
                   })
                   return
               }
               
               // 其他错误作为内部错误处理
               c.JSON(http.StatusInternalServerError, gin.H{
                   "error": "内部服务器错误",
               })
           }
       }
   }
   ```

### 3. 调试技巧

#### 3.1 日志调试
```go
// 添加调试日志
log.Printf("处理请求: %+v", req)
log.Printf("数据库查询结果: %+v", result)
log.Printf("错误信息: %v", err)
```

#### 3.2 使用 Delve 调试
```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/server/main.go

# 设置断点
break api/v1/handlers/todo.go:42

# 查看变量
print req
```

#### 3.3 性能分析
```go
import "net/http/pprof"

// 添加 pprof 路由
pprof.Register(router)

// 使用 go tool pprof 分析
go tool pprof http://localhost:8080/debug/pprof/profile
```

### 4. 常见问题解决

#### 4.1 数据库连接问题
```go
// 问题：连接池设置不当导致连接耗尽
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    ConnPool: &pool.Config{
        MaxIdleConns: 10,   // 最大空闲连接数
        MaxOpenConns: 100,  // 最大打开连接数
        ConnMaxLifetime: time.Hour, // 连接最大生命周期
    },
})
```

#### 4.2 内存泄漏问题
```go
// 问题：goroutine 泄漏
// 错误示例
go func() {
    // 没有退出机制的 goroutine
    for {
        processTask()
    }
}()

// 正确示例
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        case task := <-taskChan:
            processTask(task)
        }
    }
}()
```

#### 4.3 并发问题
   ```go
// 问题：并发访问共享资源
var counter int  // 共享资源

// 错误示例
go func() {
    counter++  // 并发不安全
}()

// 正确示例
var mu sync.Mutex
go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()
```

### 5. 最佳实践示例

#### 5.1 配置管理
```go
// 使用结构化的配置
type Config struct {
    Server struct {
        Port int    `yaml:"port"`
        Mode string `yaml:"mode"`
    } `yaml:"server"`
    
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

// 加载配置
func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("解析配置失败: %w", err)
    }
    
    return &config, nil
}
```

#### 5.2 数据库操作
```go
// 使用事务确保数据一致性
func (r *TodoRepository) CreateWithCategory(todo *Todo, category *Category) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建分类
        if err := tx.Create(category).Error; err != nil {
            return err
        }
        
        // 2. 创建待办事项
        todo.CategoryID = category.ID
        if err := tx.Create(todo).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

#### 5.3 缓存使用
```go
// 使用缓存减少数据库查询
func (s *TodoService) GetByID(id uint) (*Todo, error) {
    // 1. 尝试从缓存获取
    key := fmt.Sprintf("todo:%d", id)
    if todo, err := s.cache.Get(key); err == nil {
        return todo.(*Todo), nil
    }
    
    // 2. 缓存未命中，从数据库查询
    todo, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    s.cache.Set(key, todo, time.Hour)
    
    return todo, nil
}
```

### 6. 项目开发流程详解

#### 6.1 需求分析与设计
1. 业务需求分析
   ```
   1. 用户故事
      - 作为用户，我希望能注册/登录系统
      - 作为用户，我希望能创建待办事项
      - 作为用户，我希望能设置提醒
   
   2. 功能需求
      - 用户认证
      - 待办事项 CRUD
      - 分类管理
      - 提醒通知
   ```

2. 技术方案设计
   ```
   1. 技术选型
      - Web 框架：Gin
      - ORM：GORM
      - 缓存：Redis
      - 认证：JWT
   
   2. 架构设计
      - 分层架构
      - RESTful API
      - 微服务准备
   ```

#### 6.2 开发环境搭建

1. 项目初始化
   ```bash
   # 1. 创建项目目录
   mkdir -p todo/{cmd,api,internal,pkg}
   cd todo

   # 2. 初始化 Go 模块
   go mod init todo

   # 3. 安装基础依赖
   go get -u github.com/gin-gonic/gin
   go get -u gorm.io/gorm
   go get -u github.com/spf13/viper
   ```

2. 开发工具配置
     ```json
   // .vscode/settings.json
   {
       "go.useLanguageServer": true,
       "go.lintTool": "golangci-lint",
       "go.formatTool": "gofmt",
       "[go]": {
           "editor.formatOnSave": true,
           "editor.codeActionsOnSave": {
               "source.organizeImports": true
           }
       }
   }
   ```

#### 6.3 功能实现步骤

1. 数据模型定义
   ```go
   // internal/models/todo.go
   type Todo struct {
       gorm.Model
       Title       string    `gorm:"size:128;not null"`
       Description string    `gorm:"size:1024"`
       DueDate     time.Time `gorm:"index"`
       Completed   bool      `gorm:"default:false"`
       UserID      uint      `gorm:"index;not null"`
       CategoryID  uint      `gorm:"index"`
   }

   // 模型方法
   func (t *Todo) BeforeCreate(tx *gorm.DB) error {
       // 创建前的验证
       if strings.TrimSpace(t.Title) == "" {
           return errors.New("标题不能为空")
       }
       return nil
   }
   ```

2. 仓储层实现
   ```go
   // internal/repository/todo.go
   type TodoRepository interface {
       Create(todo *models.Todo) error
       Update(todo *models.Todo) error
       Delete(id uint) error
       FindByID(id uint) (*models.Todo, error)
       FindByUser(userID uint, page, pageSize int) ([]models.Todo, int64, error)
   }

   // internal/repository/impl/todo.go
   type todoRepository struct {
       db *gorm.DB
   }

   func (r *todoRepository) FindByUser(userID uint, page, pageSize int) ([]models.Todo, int64, error) {
       var todos []models.Todo
       var total int64
       
       // 1. 计算总数
       if err := r.db.Model(&models.Todo{}).
           Where("user_id = ?", userID).
           Count(&total).Error; err != nil {
           return nil, 0, err
       }
       
       // 2. 分页查询
       offset := (page - 1) * pageSize
       if err := r.db.Where("user_id = ?", userID).
           Offset(offset).
           Limit(pageSize).
           Find(&todos).Error; err != nil {
           return nil, 0, err
       }
       
       return todos, total, nil
   }
   ```

3. 服务层实现
   ```go
   // internal/service/todo.go
   type TodoService interface {
       Create(userID uint, req *dto.CreateTodoRequest) (*models.Todo, error)
       Update(userID uint, id uint, req *dto.UpdateTodoRequest) (*models.Todo, error)
       Delete(userID uint, id uint) error
       GetByID(userID uint, id uint) (*models.Todo, error)
       List(userID uint, page, pageSize int) (*dto.TodoListResponse, error)
   }

   // internal/service/impl/todo.go
   type todoService struct {
       repo   repository.TodoRepository
       cache  cache.Cache
       logger logger.Logger
   }

   func (s *todoService) Create(userID uint, req *dto.CreateTodoRequest) (*models.Todo, error) {
       // 1. 参数验证
       if err := req.Validate(); err != nil {
           return nil, err
       }
       
       // 2. 构建模型
       todo := &models.Todo{
           Title:       req.Title,
           Description: req.Description,
           UserID:      userID,
           DueDate:     req.DueDate,
       }
       
       // 3. 保存数据
       if err := s.repo.Create(todo); err != nil {
           s.logger.Errorf("创建待办事项失败: %v", err)
           return nil, err
       }
       
       // 4. 清除缓存
       s.cache.Delete(fmt.Sprintf("todo:list:%d", userID))
       
       return todo, nil
   }
   ```

4. 处理器实现
   ```go
   // api/v1/handlers/todo.go
   type TodoHandler struct {
       todoService service.TodoService
       logger      logger.Logger
   }

   func (h *TodoHandler) Create(c *gin.Context) {
       // 1. 获取用户信息
       userID := auth.GetUserID(c)
       
       // 2. 绑定请求数据
       var req dto.CreateTodoRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           h.logger.Errorf("绑定请求数据失败: %v", err)
           c.JSON(http.StatusBadRequest, gin.H{
               "error": "无效的请求数据",
           })
           return
       }
       
       // 3. 调用服务层
       todo, err := h.todoService.Create(userID, &req)
       if err != nil {
           h.logger.Errorf("创建待办事项失败: %v", err)
           c.JSON(http.StatusInternalServerError, gin.H{
               "error": "创建待办事项失败",
           })
           return
       }
       
       // 4. 返回结果
       c.JSON(http.StatusOK, todo)
   }
   ```

#### 6.4 测试编写指南

1. 单元测试
   ```go
   // internal/service/impl/todo_test.go
   func TestTodoService_Create(t *testing.T) {
       // 1. 准备测试数据
       ctrl := gomock.NewController(t)
       defer ctrl.Finish()
       
       mockRepo := mock_repository.NewMockTodoRepository(ctrl)
       mockCache := mock_cache.NewMockCache(ctrl)
       mockLogger := mock_logger.NewMockLogger(ctrl)
       
       service := NewTodoService(mockRepo, mockCache, mockLogger)
       
       req := &dto.CreateTodoRequest{
           Title:       "测试待办",
           Description: "这是一个测试",
       }
       
       // 2. 设置模拟行为
       mockRepo.EXPECT().
           Create(gomock.Any()).
           Return(nil)
           
       mockCache.EXPECT().
           Delete(gomock.Any()).
           Return(nil)
       
       // 3. 执行测试
       todo, err := service.Create(1, req)
       
       // 4. 验证结果
       assert.NoError(t, err)
       assert.NotNil(t, todo)
       assert.Equal(t, req.Title, todo.Title)
   }
   ```

2. 集成测试
   ```go
   // test/integration/todo_test.go
   func TestTodoAPI_Create(t *testing.T) {
       // 1. 设置测试环境
       app := setupTestApp()
       
       // 2. 创建测试客户端
       client := &http.Client{}
       
       // 3. 执行测试用例
       t.Run("创建待办事项", func(t *testing.T) {
           // 准备请求数据
           reqBody := strings.NewReader(`{
               "title": "测试待办",
               "description": "这是一个测试"
           }`)
           
           // 创建请求
           req, err := http.NewRequest(
               "POST",
               "http://localhost:8080/api/v1/todos",
               reqBody,
           )
           require.NoError(t, err)
           
           // 设置认证信息
           req.Header.Set("Authorization", "Bearer "+testToken)
           req.Header.Set("Content-Type", "application/json")
           
           // 发送请求
           resp, err := client.Do(req)
           require.NoError(t, err)
           defer resp.Body.Close()
           
           // 验证响应
           assert.Equal(t, http.StatusOK, resp.StatusCode)
           
           var todo models.Todo
           err = json.NewDecoder(resp.Body).Decode(&todo)
           require.NoError(t, err)
           
           assert.NotZero(t, todo.ID)
           assert.Equal(t, "测试待办", todo.Title)
       })
   }
   ```

### 7. 性能优化指南

#### 7.1 数据库优化

1. 索引优化
   ```sql
   -- 创建复合索引
   CREATE INDEX idx_user_completed_created 
   ON todos (user_id, completed, created_at);
   
   -- GORM 中定义索引
   type Todo struct {
       gorm.Model
       UserID    uint      `gorm:"index:idx_user_completed_created,priority:1"`
       Completed bool      `gorm:"index:idx_user_completed_created,priority:2"`
       CreatedAt time.Time `gorm:"index:idx_user_completed_created,priority:3"`
   }
   ```

2. 查询优化
   ```go
   // 使用预加载减少查询次数
   todos, err := db.Preload("Category").
       Preload("Reminders").
       Where("user_id = ?", userID).
       Find(&todos).Error

   // 使用 Select 指定需要的字段
   var todos []Todo
   db.Select("id", "title", "completed").
       Where("user_id = ?", userID).
       Find(&todos)
   ```

#### 7.2 缓存优化

1. 缓存策略
   ```go
   // 1. 缓存接口
   type Cache interface {
       Get(key string) (interface{}, error)
       Set(key string, value interface{}, expiration time.Duration) error
       Delete(key string) error
   }

   // 2. 实现缓存
   func (s *TodoService) GetByID(id uint) (*Todo, error) {
       // 缓存键
       cacheKey := fmt.Sprintf("todo:%d", id)
       
       // 1. 尝试从缓存获取
       if cached, err := s.cache.Get(cacheKey); err == nil {
           return cached.(*Todo), nil
       }
       
       // 2. 从数据库获取
       todo, err := s.repo.FindByID(id)
       if err != nil {
           return nil, err
       }
       
       // 3. 写入缓存
       s.cache.Set(cacheKey, todo, time.Hour)
       
       return todo, nil
   }
   ```

2. 缓存更新
   ```go
   // 更新待办事项时更新缓存
   func (s *TodoService) Update(id uint, req *UpdateTodoRequest) error {
       // 1. 更新数据库
       if err := s.repo.Update(id, req); err != nil {
           return err
       }
       
       // 2. 删除相关缓存
       cacheKeys := []string{
           fmt.Sprintf("todo:%d", id),
           fmt.Sprintf("todo:list:%d", req.UserID),
       }
       
       for _, key := range cacheKeys {
           s.cache.Delete(key)
       }
       
       return nil
   }
   ```

#### 7.3 并发优化

1. 连接池配置
   ```go
   // 数据库连接池
   sqlDB, err := db.DB()
   if err != nil {
       return err
   }
   
   // 设置连接池参数
   sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
   sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
   sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
   ```

2. 并发控制
   ```go
   // 使用工作池处理任务
   type WorkerPool struct {
       workers  int
       taskChan chan Task
       wg       sync.WaitGroup
   }

   func NewWorkerPool(workers int) *WorkerPool {
       return &WorkerPool{
           workers:  workers,
           taskChan: make(chan Task, workers*2),
       }
   }

   func (p *WorkerPool) Start() {
       for i := 0; i < p.workers; i++ {
           p.wg.Add(1)
           go func() {
               defer p.wg.Done()
               for task := range p.taskChan {
                   task.Process()
               }
           }()
       }
   }
   ```

### 8. 监控与日志

#### 8.1 性能监控
```go
// 添加 Prometheus 指标
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "path"},
    )
)

// 使用中间件记录指标
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
        ).Observe(duration)
    }
}
```

#### 8.2 日志记录
```go
// 结构化日志
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
    With(fields ...Field) Logger
}

// 实现日志记录
func (s *TodoService) Create(req *CreateTodoRequest) error {
    logger := s.logger.With(
        log.String("action", "create_todo"),
        log.String("user_id", req.UserID),
    )
    
    logger.Info("开始创建待办事项")
    
    // 业务逻辑...
    
    if err != nil {
        logger.Error("创建待办事项失败", 
            log.Error(err),
            log.String("title", req.Title),
        )
        return err
    }
    
    logger.Info("待办事项创建成功", 
        log.Int("todo_id", todo.ID),
    )
    return nil
}
```

## 9. API 开发与测试

### 9.1 API 文档规范

1. Swagger 注解规范
```go
// @Summary 创建待办事项
// @Description 创建一个新的待办事项
// @Tags todos
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param todo body CreateTodoRequest true "待办事项信息"
// @Success 200 {object} Todo
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context)
```

2. API 响应格式
```go
// 成功响应
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 错误响应
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
```

### 9.2 API 测试方法

1. curl 测试示例
```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{
"username": "test",
"password": "test123",
"email": "test@example.com"
}'
```

```bash
# 登录获取token
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
"username": "test",
"password": "test123",
}'
```

2. 单元测试
```go
func TestTodoAPI_Create(t *testing.T) {
    // 1. 设置测试环境
    router := setupTestRouter()
    
    // 2. 创建测试请求
    body := `{
        "title": "测试待办",
        "description": "这是一个测试"
    }`
    
    req := httptest.NewRequest("POST", "/api/v1/todos", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+testToken)
    
    // 3. 记录响应
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // 4. 验证结果
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response Response
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "创建成功", response.Message)
}
```

1. 集成测试
```go
func TestTodoAPI_Integration(t *testing.T) {
    // 1. 启动测试服务器
    ts := httptest.NewServer(setupTestRouter())
    defer ts.Close()
    
    // 2. 创建 HTTP 客户端
    client := &http.Client{}
    
    // 3. 测试创建待办事项
    t.Run("Create Todo", func(t *testing.T) {
        // 发送请求
        resp, err := client.Post(
            ts.URL+"/api/v1/todos",
            "application/json",
            strings.NewReader(`{"title":"测试待办"}`),
        )
        
        // 验证结果
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })
}
```

### 9.3 API 性能测试

1. 基准测试
```go
func BenchmarkTodoAPI_Create(b *testing.B) {
    router := setupTestRouter()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        req := httptest.NewRequest("POST", "/api/v1/todos",
            strings.NewReader(`{"title":"测试待办"}`))
        
        router.ServeHTTP(w, req)
    }
}
```

2. 负载测试
```go
func TestTodoAPI_Load(t *testing.T) {
    // 1. 设置并发用户数
    users := 100
    
    // 2. 创建等待组
    var wg sync.WaitGroup
    wg.Add(users)
    
    // 3. 模拟并发请求
    for i := 0; i < users; i++ {
        go func() {
            defer wg.Done()
            
            // 发送请求
            resp, err := http.Post(
                "http://localhost:8080/api/v1/todos",
                "application/json",
                strings.NewReader(`{"title":"测试待办"}`),
            )
            
            // 验证响应
            assert.NoError(t, err)
            assert.Equal(t, http.StatusOK, resp.StatusCode)
        }()
    }
    
    wg.Wait()
}
```

### 9.4 API 安全测试

1. 认证测试
```go
func TestTodoAPI_Auth(t *testing.T) {
    router := setupTestRouter()
    
    // 1. 测试无 token 访问
    t.Run("No Token", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/todos", nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })
    
    // 2. 测试无效 token
    t.Run("Invalid Token", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/todos", nil)
        req.Header.Set("Authorization", "Bearer invalid-token")
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })
}
```

2. 权限测试
```go
func TestTodoAPI_Authorization(t *testing.T) {
    router := setupTestRouter()
    
    // 测试越权访问
    t.Run("Unauthorized Access", func(t *testing.T) {
        // 创建测试用户和待办事项
        user1Token := createTestUser(t, "user1")
        user2Token := createTestUser(t, "user2")
        todoID := createTestTodo(t, user1Token)
        
        // 尝试使用 user2 的 token 访问 user1 的待办事项
        req := httptest.NewRequest("GET", "/api/v1/todos/"+todoID, nil)
        req.Header.Set("Authorization", "Bearer "+user2Token)
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusForbidden, w.Code)
    })
}
```

## 10. 高级特性实现指南

### 10.1 分布式锁实现
```go
// pkg/lock/distributed_lock.go

// 分布式锁接口
type DistributedLock interface {
    Lock(ctx context.Context, key string) error
    Unlock(ctx context.Context, key string) error
    TryLock(ctx context.Context, key string) (bool, error)
}

// Redis 实现
type RedisLock struct {
    client *redis.Client
    ttl    time.Duration
}

func (l *RedisLock) Lock(ctx context.Context, key string) error {
    // 使用 SET NX 实现加锁
    success, err := l.client.SetNX(ctx, 
        "lock:"+key, 
        uuid.New().String(), 
        l.ttl,
    ).Result()
    
    if err != nil {
        return fmt.Errorf("获取锁失败: %w", err)
    }
    
    if !success {
        return ErrLockAlreadyHeld
    }
    
    return nil
}

// 使用示例
func (s *TodoService) UpdateWithLock(id uint, req UpdateRequest) error {
    ctx := context.Background()
    key := fmt.Sprintf("todo:%d", id)
    
    // 获取锁
    if err := s.lock.Lock(ctx, key); err != nil {
        return fmt.Errorf("获取锁失败: %w", err)
    }
    defer s.lock.Unlock(ctx, key)
    
    // 执行更新操作
    return s.repo.Update(id, req)
}
```

### 10.2 消息队列集成
```go
// pkg/queue/task_queue.go

// 任务队列接口
type TaskQueue interface {
    Push(ctx context.Context, task Task) error
    Pop(ctx context.Context) (*Task, error)
    Process(ctx context.Context, handler TaskHandler) error
}

// Redis 实现
type RedisQueue struct {
    client *redis.Client
    key    string
}

func (q *RedisQueue) Push(ctx context.Context, task Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    return q.client.LPush(ctx, q.key, data).Err()
}

// 使用示例
func (s *TodoService) CreateAsync(req CreateRequest) error {
    task := Task{
        Type: TaskTypeTodoCreate,
        Data: req,
    }
    
    return s.queue.Push(context.Background(), task)
}
```

### 10.3 缓存穿透防护
```go
// pkg/cache/bloom_filter.go

// 布隆过滤器
type BloomFilter struct {
    bits    *bitset.BitSet
    hashFns []hash.Hash64
}

func (bf *BloomFilter) Add(key string) {
    for _, h := range bf.hashFns {
        h.Write([]byte(key))
        i := h.Sum64() % uint64(bf.bits.Len())
        bf.bits.Set(uint(i))
        h.Reset()
    }
}

// 缓存服务
type CacheService struct {
    cache  Cache
    bloom  *BloomFilter
    repo   Repository
}

func (s *CacheService) Get(key string) (interface{}, error) {
    // 1. 检查布隆过滤器
    if !s.bloom.MightContain(key) {
        return nil, ErrNotFound
    }
    
    // 2. 查询缓存
    if val, err := s.cache.Get(key); err == nil {
        return val, nil
    }
    
    // 3. 查询数据库
    val, err := s.repo.Find(key)
    if err != nil {
        return nil, err
    }
    
    // 4. 写入缓存
    s.cache.Set(key, val, time.Hour)
    return val, nil
}
```

### 10.4 限流器实现
```go
// pkg/middleware/rate_limit.go

// 令牌桶限流器
type TokenBucket struct {
    rate     float64
    capacity int64
    tokens   int64
    lastTime time.Time
    mu       sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    now := time.Now()
    
    // 计算新增令牌
    elapsed := now.Sub(tb.lastTime).Seconds()
    newTokens := int64(elapsed * tb.rate)
    
    // 更新令牌数
    if tb.tokens+newTokens > tb.capacity {
        tb.tokens = tb.capacity
    } else {
        tb.tokens += newTokens
    }
    
    tb.lastTime = now
    
    // 判断是否允许请求
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    
    return false
}

// 中间件实现
func RateLimitMiddleware(limiter *TokenBucket) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatus(http.StatusTooManyRequests)
            return
        }
        c.Next()
    }
}
```

### 10.5 链路追踪实现
```go
// pkg/trace/tracer.go

// 追踪上下文
type TraceContext struct {
    TraceID    string
    SpanID     string
    ParentID   string
    StartTime  time.Time
    EndTime    time.Time
    Tags       map[string]string
}

// 追踪器
type Tracer interface {
    StartSpan(name string) *TraceContext
    FinishSpan(ctx *TraceContext)
}

// 中间件实现
func TracingMiddleware(tracer Tracer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 创建追踪上下文
        span := tracer.StartSpan(c.Request.URL.Path)
        defer tracer.FinishSpan(span)
        
        // 设置追踪信息
        c.Set("trace_context", span)
        
        c.Next()
    }
}

// 服务层使用
func (s *TodoService) Create(req CreateRequest) error {
    // 获取追踪上下文
    span := s.tracer.StartSpan("TodoService.Create")
    defer s.tracer.FinishSpan(span)
    
    // 添加标签
    span.Tags["user_id"] = req.UserID
    span.Tags["todo_title"] = req.Title
    
    // 执行业务逻辑
    return s.repo.Create(req)
}
```

### 10.6 优雅关闭实现
```go
// pkg/server/graceful.go

type Server struct {
    http   *http.Server
    done   chan struct{}
    logger Logger
}

func (s *Server) Start() error {
    // 启动 HTTP 服务器
    go func() {
        if err := s.http.ListenAndServe(); err != nil && 
           err != http.ErrServerClosed {
            s.logger.Errorf("HTTP server error: %v", err)
        }
    }()
    
    // 监听信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    <-quit
    s.logger.Info("Shutting down server...")
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := s.http.Shutdown(ctx); err != nil {
        s.logger.Errorf("Server forced to shutdown: %v", err)
        return err
    }
    
    s.logger.Info("Server exiting")
    close(s.done)
    return nil
}
```

### 10.7 熔断器实现
```go
// pkg/circuit/breaker.go

type CircuitBreaker struct {
    failures  int64
    threshold int64
    timeout   time.Duration
    lastErr   error
    mu        sync.RWMutex
    state     State
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.AllowRequest() {
        return ErrCircuitOpen
    }
    
    err := fn()
    cb.RecordResult(err)
    return err
}

// 服务层使用
func (s *TodoService) CreateWithBreaker(req CreateRequest) error {
    return s.breaker.Execute(func() error {
        return s.repo.Create(req)
    })
}
```

## 11. 项目最佳实践总结

### 11.1 代码组织
1. 清晰的目录结构
2. 合理的包划分
3. 统一的命名规范
4. 完善的注释文档

### 11.2 错误处理
1. 统一错误定义
2. 错误包装传递
3. 合理的错误返回
4. 日志记录关联

### 11.3 性能优化
1. 合理使用缓存
2. 数据库优化
3. 并发控制
4. 资源池化

### 11.4 安全考虑
1. 输入验证
2. 权限控制
3. 敏感信息加密
4. 安全审计日志

### 11.5 可维护性
1. 模块化设计
2. 依赖注入
3. 接口抽象
4. 测试覆盖

### 11.6 可扩展性
1. 水平扩展支持
2. 配置外部化
3. 插件化设计
4. 版本兼容性