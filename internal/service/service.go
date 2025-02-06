package service

import (
	"context"
	"todo/api/v1/dto/category"
	"todo/api/v1/dto/reminder"
	"todo/api/v1/dto/todo"
	"todo/internal/models"
	"todo/internal/repository"
	"todo/internal/service/impl"
	"todo/pkg/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewAuthService 创建新的认证服务实例
func NewAuthService(db *gorm.DB, rdb *redis.Client, jwtCfg *config.JWTConfig) AuthService {
	userRepo := repository.NewUserRepository(db)
	return impl.NewAuthService(userRepo, jwtCfg)
}

// NewTodoService 创建新的待办事项服务实例
func NewTodoService(db *gorm.DB) TodoService {
	todoRepo := repository.NewTodoRepository(db)
	return impl.NewTodoService(todoRepo)
}

// NewCategoryService 创建新的分类服务实例
func NewCategoryService(db *gorm.DB) CategoryService {
	categoryRepo := repository.NewCategoryRepository(db)
	svc := impl.NewCategoryService(categoryRepo)
	return &categoryServiceWrapper{svc}
}

// NewReminderService 创建新的提醒服务实例
func NewReminderService(db *gorm.DB) ReminderService {
	reminderRepo := repository.NewReminderRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	svc := impl.NewReminderService(reminderRepo, todoRepo)
	return &reminderServiceWrapper{svc}
}

// Wrapper types
type todoServiceWrapper struct {
	svc *impl.TodoService
}

type categoryServiceWrapper struct {
	svc *impl.CategoryService
}

type reminderServiceWrapper struct {
	svc *impl.ReminderService
}

// 服务层包装器实现
// 这里使用包装器模式来统一处理DTO与内部模型的转换

// TodoService wrapper implementations
func (w *todoServiceWrapper) Create(ctx context.Context, userID uint, req *todo.CreateRequest) (uint, error) {
	return w.svc.Create(ctx, userID, req)
}

func (w *todoServiceWrapper) List(ctx context.Context, userID uint) ([]*models.Todo, error) {
	return w.svc.List(ctx, userID)
}

func (w *todoServiceWrapper) Get(ctx context.Context, id, userID uint) (*models.Todo, error) {
	return w.svc.Get(ctx, id, userID)
}

func (w *todoServiceWrapper) Update(ctx context.Context, id, userID uint, req *todo.UpdateRequest) error {
	return w.svc.Update(ctx, id, userID, req)
}

func (w *todoServiceWrapper) Delete(ctx context.Context, id, userID uint) error {
	return w.svc.Delete(ctx, id, userID)
}

// CategoryService wrapper implementations
func (w *categoryServiceWrapper) Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error) {
	return w.svc.Create(ctx, userID, req)
}

func (w *categoryServiceWrapper) List(ctx context.Context, userID uint) ([]*models.Category, error) {
	return w.svc.List(ctx, userID)
}

func (w *categoryServiceWrapper) Get(ctx context.Context, id, userID uint) (*models.Category, error) {
	return w.svc.Get(ctx, id, userID)
}

func (w *categoryServiceWrapper) Update(ctx context.Context, id, userID uint, req *category.UpdateRequest) error {
	return w.svc.Update(ctx, userID, id, req)
}

func (w *categoryServiceWrapper) Delete(ctx context.Context, id, userID uint) error {
	return w.svc.Delete(ctx, userID, id)
}

// ReminderService wrapper implementations
func (w *reminderServiceWrapper) Create(ctx context.Context, userID uint, req *reminder.CreateRequest) (uint, error) {
	return w.svc.Create(ctx, userID, req)
}

func (w *reminderServiceWrapper) ListByTodoID(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	return w.svc.ListByTodoID(ctx, todoID)
}

func (w *reminderServiceWrapper) Get(ctx context.Context, id, userID uint) (*models.Reminder, error) {
	return w.svc.Get(ctx, id, userID)
}

func (w *reminderServiceWrapper) Update(ctx context.Context, id, userID uint, req *reminder.UpdateRequest) error {
	return w.svc.Update(ctx, id, userID, req)
}

func (w *reminderServiceWrapper) Delete(ctx context.Context, id, userID uint) error {
	return w.svc.Delete(ctx, id, userID)
}
