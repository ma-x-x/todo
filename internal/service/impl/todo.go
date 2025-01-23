package impl

import (
	"context"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
)

type TodoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

func (s *TodoService) Create(ctx context.Context, userID uint, req *todo.CreateRequest) (uint, error) {
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    models.Priority(req.Priority),
		UserID:      userID,
		CategoryID:  req.CategoryID,
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return 0, err
	}

	return todo.ID, nil
}

func (s *TodoService) GetByID(ctx context.Context, id, userID uint) (*models.Todo, error) {
	return s.todoRepo.GetByID(ctx, id)
}

func (s *TodoService) GetTodoRepo() repository.TodoRepository {
	return s.todoRepo
}
