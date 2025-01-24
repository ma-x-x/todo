package impl

import (
	"context"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
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
	todo := &models.Todo{
		Title:       req.Title,                     // 待办事项标题
		Description: req.Description,               // 待办事项描述
		Priority:    models.Priority(req.Priority), // 优先级
		UserID:      userID,                        // 所属用户ID
		CategoryID:  req.CategoryID,                // 所属分类ID
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return 0, err
	}

	return todo.ID, nil
}

// GetByID 根据ID获取待办事项
//
// Parameters:
//   - ctx: 上下文信息
//   - id: 待办事项ID
//   - userID: 用户ID
//
// Returns:
//   - *models.Todo: 返回待办事项信息
//   - error: 可能的错误信息
func (s *TodoService) GetByID(ctx context.Context, id, userID uint) (*models.Todo, error) {
	return s.todoRepo.GetByID(ctx, id)
}

// GetTodoRepo 获取待办事项仓库实例
//
// Returns:
//   - repository.TodoRepository: 返回待办事项仓库接口
func (s *TodoService) GetTodoRepo() repository.TodoRepository {
	return s.todoRepo
}
