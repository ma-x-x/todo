package category

import "todo/internal/models"

// ListResponse 分类列表响应
type ListResponse struct {
	Total int64              `json:"total"`     // 总数
	Items []*models.Category `json:"items"`     // 分类列表
}
