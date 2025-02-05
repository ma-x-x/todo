package category

// CreateRequest 创建分类请求
type CreateRequest struct {
	Name  string `json:"name" binding:"required,max=32"`
	Color string `json:"color" binding:"omitempty,max=7"`
}

// CreateResponse 创建分类响应
type CreateResponse struct {
	ID uint `json:"id"`
}
