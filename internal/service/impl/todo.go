package impl

import (
	"context"
	"time"
	"todo/api/v1/dto/todo"
	"todo/internal/models"
	"todo/internal/repository/interfaces"
	"todo/pkg/errors"
)

// TodoService 待办事项服务结构体
// 负责处理所有与待办事项相关的业务逻辑
type TodoService struct {
	todoRepo interfaces.TodoRepository
}

// NewTodoService 创建一个新的待办事项服务实例
//
// Parameters:
//   - repo: 待办事项仓库实现
//
// Returns:
//   - *TodoService: 返回待办事项服务实例
func NewTodoService(repo interfaces.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: repo,
	}
}

// Create 创建新的待办事项
//
// Parameters:
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - req: 创建待办事项的请求数据
//
// Returns:
//   - uint: 返回新创建的待办事项ID
//   - error: 可能的错误信息
func (s *TodoService) Create(ctx context.Context, userID uint, req *todo.CreateRequest) (uint, error) {
	todoItem := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
		CategoryID:  req.CategoryID,
	}

	if req.Priority != "" {
		todoItem.Priority = models.Priority(req.Priority)
	} else {
		todoItem.Priority = models.PriorityMedium // 默认中优先级
	}

	// 只在提供了截止时间时设置
	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.DueDate)
		if err != nil {
			return 0, errors.New("无效的截止时间格式")
		}
		todoItem.DueDate = &dueDate // 使用指针
	}

	if err := s.todoRepo.Create(ctx, todoItem); err != nil {
		return 0, err
	}

	return todoItem.ID, nil
}

// List 获取用户的所有待办事项
func (s *TodoService) List(ctx context.Context, userID uint) ([]*models.Todo, error) {
	return s.todoRepo.ListByUserID(ctx, userID)
}

// Get 获取单个待办事项详情
func (s *TodoService) Get(ctx context.Context, id, userID uint) (*models.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 验证所有权
	if todo.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return todo, nil
}

// Update 更新待办事项
func (s *TodoService) Update(ctx context.Context, id, userID uint, req *todo.UpdateRequest) error {
	todoItem, err := s.todoRepo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	if req.Title != "" {
		todoItem.Title = req.Title
	}

	if req.Description != "" {
		todoItem.Description = req.Description
	}

	if req.Status != "" {
		status, err := models.ParseTodoStatus(req.Status)
		if err != nil {
			return err
		}
		todoItem.Status = status
	}

	if req.Priority != "" {
		todoItem.Priority = models.Priority(req.Priority)
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.DueDate)
		if err != nil {
			return errors.New("无效的截止时间格式")
		}
		todoItem.DueDate = &dueDate // 修改这里，使用指针
	}

	if req.CategoryID != nil {
		todoItem.CategoryID = req.CategoryID
	}

	if req.Completed != nil {
		todoItem.Completed = *req.Completed
	}

	return s.todoRepo.Update(ctx, todoItem)
}

// Delete 删除待办事项
func (s *TodoService) Delete(ctx context.Context, id, userID uint) error {
	todo, err := s.Get(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.todoRepo.Delete(ctx, todo.ID)
}

// GetTodoRepo 获取待办事项仓库实例
//
// Returns:
//   - interfaces.TodoRepository: 返回待办事项仓库接口
func (s *TodoService) GetTodoRepo() interfaces.TodoRepository {
	return s.todoRepo
}
