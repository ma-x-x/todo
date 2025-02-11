package service

import (
	"context"
	"sync"
	"todo/internal/repository/interfaces"
	"todo/internal/service/impl"
	"todo/pkg/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewAuthService 创建新的认证服务实例
func NewAuthService(userRepo interfaces.UserRepository, authRepo interfaces.AuthRepository, jwtCfg *config.JWTConfig) AuthService {
	return impl.NewAuthService(userRepo, authRepo, jwtCfg)
}

// NewTodoService 创建新的待办事项服务实例
func NewTodoService(repo interfaces.TodoRepository) TodoService {
	return impl.NewTodoService(repo)
}

// NewCategoryService 创建新的分类服务实例
func NewCategoryService(repo interfaces.CategoryRepository) CategoryService {
	return impl.NewCategoryService(repo)
}

// NewReminderService 创建新的提醒服务实例
func NewReminderService(reminderRepo interfaces.ReminderRepository, todoRepo interfaces.TodoRepository) ReminderService {
	return impl.NewReminderService(reminderRepo, todoRepo)
}

// 1. 添加服务生命周期管理
type Service interface {
	Start() error
	Stop() error
	Health() bool
}

// 2. 添加服务依赖注入容器
type Container struct {
	services map[string]Service
	mu       sync.RWMutex
}

// 3. 实现服务注册机制
func (c *Container) Register(name string, service Service) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// ServiceCollection 所有服务的集合
type ServiceCollection struct {
	Auth     AuthService
	Todo     TodoService
	Category CategoryService
	Reminder ReminderService
	db       *gorm.DB      // 添加数据库连接
	rdb      *redis.Client // 添加Redis连接
}

// CheckDatabaseHealth 检查数据库健康状态
func (s *ServiceCollection) CheckDatabaseHealth() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// CheckRedisHealth 检查Redis健康状态
func (s *ServiceCollection) CheckRedisHealth() error {
	return s.rdb.Ping(context.Background()).Err()
}

// NewServiceCollection 创建新的服务集合实例
func NewServiceCollection(
	auth AuthService,
	todo TodoService,
	category CategoryService,
	reminder ReminderService,
	db *gorm.DB,
	rdb *redis.Client,
) *ServiceCollection {
	return &ServiceCollection{
		Auth:     auth,
		Todo:     todo,
		Category: category,
		Reminder: reminder,
		db:       db,
		rdb:      rdb,
	}
}
