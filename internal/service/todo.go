package service

import (
	"context"
	"todo-demo/internal/models"
)

type TodoService interface {
	Create(ctx context.Context, userID uint, title, description string,
		priority models.Priority, categoryID *uint) (*models.Todo, error)
	List(ctx context.Context, userID uint) ([]*models.Todo, error)
	Get(ctx context.Context, id, userID uint) (*models.Todo, error)
	Update(ctx context.Context, id, userID uint, req interface{}) error
	Delete(ctx context.Context, id, userID uint) error
}
