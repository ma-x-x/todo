package category

import "todo-demo/internal/models"

type ListResponse struct {
	Total int64              `json:"total"`
	Items []*models.Category `json:"items"`
}
