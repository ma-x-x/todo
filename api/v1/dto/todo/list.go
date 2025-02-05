package todo

import "todo-demo/internal/models"

// ListResponse 待办事项列表响应
type ListResponse struct {
	// 总记录数
	// 用于前端分页显示
	Total int64 `json:"total"`

	// 待办事项列表
	// 包含当前页的所有待办事项详细信息
	Items []*models.Todo `json:"items"`
}
