package todo

import "todo/internal/models"

// DetailResponse 待办事项详情响应
type DetailResponse struct {
    *models.Todo
} 