package service

import (
	"context"
	"todo-demo/internal/models"
)

type CategoryService interface {
	Create(ctx context.Context, userID uint, name, color string) (*models.Category, error)
	List(ctx context.Context, userID uint) ([]*models.Category, error)
	Get(ctx context.Context, id, userID uint) (*models.Category, error)
	Update(ctx context.Context, id, userID uint, name, color string) error
	Delete(ctx context.Context, id, userID uint) error
}
