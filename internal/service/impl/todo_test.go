package impl

import (
	"context"
	"testing"
	"todo/api/v1/dto/todo"
	"todo/internal/models"
	"todo/pkg/errors"
)

// mockTodoRepo 模拟待办事项仓储接口
// 用于单元测试，避免依赖真实数据库
type mockTodoRepo struct {
	todos map[uint]*models.Todo // 存储待办事项的内存映射
	seq   uint                  // 自增ID序列
}

// newMockTodoRepo 创建一个新的待办事项仓储mock对象
func newMockTodoRepo() *mockTodoRepo {
	return &mockTodoRepo{
		todos: make(map[uint]*models.Todo),
		seq:   1,
	}
}

// Create 创建待办事项
func (m *mockTodoRepo) Create(ctx context.Context, todo *models.Todo) error {
	todo.ID = m.seq
	m.todos[todo.ID] = todo
	m.seq++
	return nil
}

// GetByID 根据ID获取待办事项
func (m *mockTodoRepo) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.ErrTodoNotFound
	}
	return todo, nil
}

// Delete 删除待办事项
func (m *mockTodoRepo) Delete(ctx context.Context, id uint) error {
	if _, exists := m.todos[id]; !exists {
		return errors.ErrTodoNotFound
	}
	delete(m.todos, id)
	return nil
}

// ListByUserID 获取用户的待办事项列表
func (m *mockTodoRepo) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*models.Todo, int64, error) {
	var todos []*models.Todo
	var total int64

	// 筛选出属于指定用户的待办事项
	for _, todo := range m.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	total = int64(len(todos))

	// 实现分页逻辑
	start := (page - 1) * pageSize
	end := start + pageSize
	if start < len(todos) {
		if end > len(todos) {
			end = len(todos)
		}
		todos = todos[start:end]
	} else {
		todos = []*models.Todo{}
	}
	return todos, total, nil
}

// Update 更新待办事项
func (m *mockTodoRepo) Update(ctx context.Context, todo *models.Todo) error {
	if _, exists := m.todos[todo.ID]; !exists {
		return errors.ErrTodoNotFound
	}
	m.todos[todo.ID] = todo
	return nil
}

// TestTodoService_Create 测试创建待办事项功能
func TestTodoService_Create(t *testing.T) {
	// 初始化测试环境
	todoRepo := newMockTodoRepo()
	todoService := NewTodoService(todoRepo)

	// 定义测试用例
	tests := []struct {
		name    string              // 测试用例名称
		userID  uint                // 用户ID
		req     *todo.CreateRequest // 创建请求
		wantErr error               // 期望的错误
	}{
		{
			name:   "成功创建待办事项",
			userID: 1,
			req: &todo.CreateRequest{
				Title:       "测试待办事项",
				Description: "测试描述",
				Priority:    1,
			},
			wantErr: nil,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := todoService.Create(context.Background(), tt.userID, tt.req)
			if err != tt.wantErr {
				t.Errorf("Create() 错误 = %v, 期望错误 %v", err, tt.wantErr)
			}
			if err == nil && id == 0 {
				t.Error("Create() 返回了无效的ID")
			}
		})
	}
}
