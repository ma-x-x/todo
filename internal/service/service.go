package service

import (
	"sync"
	"todo/internal/repository/interfaces"
	"todo/internal/service/impl"
	"todo/pkg/config"
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

// Services 所有服务的集合
type Services struct {
	Auth     AuthService
	Todo     TodoService
	Category CategoryService
	Reminder ReminderService
}

// NewServices 创建服务集合的实例
func NewServices(repos *interfaces.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:     NewAuthService(repos.User, repos.Auth, &cfg.JWT),
		Todo:     NewTodoService(repos.Todo),
		Category: NewCategoryService(repos.Category),
		Reminder: NewReminderService(repos.Reminder, repos.Todo),
	}
}
