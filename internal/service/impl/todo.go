package impl

import (
	"context"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/errors"
)

// TodoService 待办事项服务结构体
// 负责处理所有与待办事项相关的业务逻辑
type TodoService struct {
	todoRepo repository.TodoRepository // 待办事项数据仓库接口
}

// NewTodoService 创建一个新的待办事项服务实例
//
// Parameters:
//   - todoRepo: 待办事项仓库实现
//
// Returns:
//   - *TodoService: 返回待办事项服务实例
func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
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

	if err := s.todoRepo.Create(ctx, todoItem); err != nil {
		return 0, err
	}

	return todoItem.ID, nil
}

// List 获取用户的所有待办事项
func (s *TodoService) List(ctx context.Context, userID uint) ([]*models.Todo, error) {
	// 默认分页参数
	page := 1
	pageSize := 100
	todos, _, err := s.todoRepo.ListByUserID(ctx, userID, page, pageSize)
	return todos, err
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
	todoItem, err := s.Get(ctx, id, userID)
	if err != nil {
		return err
	}

	if req.Title != nil {
		todoItem.Title = *req.Title
	}
	if req.Description != nil {
		todoItem.Description = *req.Description
	}
	if req.Completed != nil {
		todoItem.Completed = *req.Completed
	}
	if req.Priority != nil {
		todoItem.Priority = models.Priority(*req.Priority)
	}
	if req.CategoryID != nil {
		todoItem.CategoryID = req.CategoryID
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
//   - repository.TodoRepository: 返回待办事项仓库接口
func (s *TodoService) GetTodoRepo() repository.TodoRepository {
	return s.todoRepo
}
