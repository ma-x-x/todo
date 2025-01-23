package category

import "todo-demo/internal/models"

// CreateRequest 创建分类请求
type CreateRequest struct {
	// Name 分类名称，必填，最大长度32个字符
	Name string `json:"name" binding:"required,max=32"`
	// Color 分类颜色，可选，使用十六进制颜色码，如 #FF0000
	Color string `json:"color" binding:"omitempty,max=7"`
}

// CreateResponse 创建分类响应
type CreateResponse struct {
	// Category 创建成功后返回的分类信息
	Category *models.Category `json:"category"`
}
