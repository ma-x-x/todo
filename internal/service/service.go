package service

import (
	"context"
	"errors"
	"time"
	"todo-demo/api/v1/dto/category"
	"todo-demo/api/v1/dto/reminder"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/internal/service/impl"
	"todo-demo/pkg/config"

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
	svc := impl.NewTodoService(todoRepo)
	return &todoServiceWrapper{svc}
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
func (w *todoServiceWrapper) Create(ctx context.Context, userID uint, title, description string,
	priority models.Priority, categoryID *uint) (*models.Todo, error) {
	// 将输入参数转换为DTO请求对象
	req := &todo.CreateRequest{
		Title:       title,         // 待办事项标题
		Description: description,   // 待办事项描述
		Priority:    int(priority), // 优先级(0:低, 1:中, 2:高)
		CategoryID:  categoryID,    // 分类ID(可选)
	}

	// 调用实际的服务层创建方法
	id, err := w.svc.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// 查询并返回新创建的待办事项
	return w.svc.GetByID(ctx, id, userID)
}

func (w *todoServiceWrapper) List(ctx context.Context, userID uint) ([]*models.Todo, error) {
	todos, _, err := w.svc.GetTodoRepo().ListByUserID(ctx, userID, 1, 20)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (w *todoServiceWrapper) Get(ctx context.Context, id, userID uint) (*models.Todo, error) {
	return w.svc.GetByID(ctx, id, userID)
}

func (w *todoServiceWrapper) Update(ctx context.Context, id, userID uint, req interface{}) error {
	// 先获取待办事项
	todoItem, err := w.svc.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// 将通用请求转换为具体的更新请求类型
	updateReq, ok := req.(*todo.UpdateRequest)
	if !ok {
		return errors.New("invalid update request type")
	}

	// 按需更新字段
	// 只有请求中包含的字段才会被更新
	if updateReq.Title != nil {
		todoItem.Title = *updateReq.Title
	}
	if updateReq.Description != nil {
		todoItem.Description = *updateReq.Description
	}
	if updateReq.Completed != nil {
		todoItem.Completed = *updateReq.Completed
	}
	if updateReq.CategoryID != nil {
		todoItem.CategoryID = updateReq.CategoryID
	}
	if updateReq.Priority != nil {
		todoItem.Priority = models.Priority(*updateReq.Priority)
	}

	// 保存更新后的待办事项
	return w.svc.GetTodoRepo().Update(ctx, todoItem)
}

func (w *todoServiceWrapper) Delete(ctx context.Context, id, userID uint) error {
	todo, err := w.svc.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return w.svc.GetTodoRepo().Delete(ctx, todo.ID)
}

// CategoryService wrapper implementations
func (w *categoryServiceWrapper) Create(ctx context.Context, userID uint, name, color string) (*models.Category, error) {
	id, err := w.svc.Create(ctx, userID, &category.CreateRequest{
		Name:  name,
		Color: color,
	})
	if err != nil {
		return nil, err
	}
	return w.svc.GetByID(ctx, userID, id)
}

func (w *categoryServiceWrapper) List(ctx context.Context, userID uint) ([]*models.Category, error) {
	return w.svc.List(ctx, userID)
}

func (w *categoryServiceWrapper) Get(ctx context.Context, id, userID uint) (*models.Category, error) {
	return w.svc.GetByID(ctx, userID, id)
}

func (w *categoryServiceWrapper) Update(ctx context.Context, id, userID uint, name, color string) error {
	return w.svc.Update(ctx, userID, id, &category.UpdateRequest{
		Name:  &name,
		Color: &color,
	})
}

func (w *categoryServiceWrapper) Delete(ctx context.Context, id, userID uint) error {
	return w.svc.Delete(ctx, userID, id)
}

// ReminderService wrapper implementations
func (w *reminderServiceWrapper) Create(ctx context.Context, userID, todoID uint, remindAt time.Time,
	remindType models.ReminderType, notifyType models.NotifyType) (*models.Reminder, error) {
	// 将输入参数转换为提醒创建请求
	req := &reminder.CreateRequest{
		TodoID:     todoID,              // 关联的待办事项ID
		RemindAt:   remindAt,            // 提醒时间
		RemindType: remindType.String(), // 提醒类型(once/daily/weekly)
		NotifyType: notifyType.String(), // 通知方式(email/push)
	}

	// 创建提醒并返回完整的提醒信息
	id, err := w.svc.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	return w.svc.GetByID(ctx, userID, id)
}

func (w *reminderServiceWrapper) List(ctx context.Context, todoID uint) ([]*models.Reminder, error) {
	// 这里我们假设用户ID的检查在服务层处理
	return w.svc.ListByTodoID(ctx, 0, todoID)
}

func (w *reminderServiceWrapper) Get(ctx context.Context, id uint) (*models.Reminder, error) {
	return w.svc.GetByID(ctx, 0, id)
}

func (w *reminderServiceWrapper) Update(ctx context.Context, id uint, remindAt time.Time,
	remindType models.ReminderType, notifyType models.NotifyType) error {
	// 构建更新请求
	req := &reminder.UpdateRequest{
		RemindAt:   remindAt,            // 新的提醒时间
		RemindType: remindType.String(), // 新的提醒类型
		NotifyType: notifyType.String(), // 新的通知方式
	}

	// 验证更新请求的合法性
	if err := req.Validate(); err != nil {
		return err
	}

	// 获取并更新提醒信息
	reminder, err := w.svc.GetByID(ctx, 0, id)
	if err != nil {
		return err
	}

	reminder.RemindAt = remindAt
	reminder.RemindType = remindType
	reminder.NotifyType = notifyType
	return w.svc.GetReminderRepo().Update(ctx, reminder)
}

func (w *reminderServiceWrapper) Delete(ctx context.Context, id uint) error {
	reminder, err := w.svc.GetByID(ctx, 0, id)
	if err != nil {
		return err
	}
	return w.svc.GetReminderRepo().Delete(ctx, reminder.ID)
}
