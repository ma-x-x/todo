package impl

import (
	"context"
	"todo-demo/api/v1/dto/category"
	"todo-demo/internal/models"
	"todo-demo/internal/repository"
	"todo-demo/pkg/errors"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(ctx context.Context, userID uint, req *category.CreateRequest) (uint, error) {
	category := &models.Category{
		Name:   req.Name,
		Color:  req.Color,
		UserID: userID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (s *CategoryService) GetByID(ctx context.Context, userID, categoryID uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return category, nil
}

func (s *CategoryService) List(ctx context.Context, userID uint) ([]*models.Category, error) {
	return s.categoryRepo.ListByUserID(ctx, userID)
}

func (s *CategoryService) Update(ctx context.Context, userID, categoryID uint, req *category.UpdateRequest) error {
	category, err := s.GetByID(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Color != nil {
		category.Color = *req.Color
	}

	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, userID, categoryID uint) error {
	category, err := s.GetByID(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, category.ID)
}
