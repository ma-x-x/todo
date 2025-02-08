package category

// CreateRequest 创建分类请求
// @Description 创建分类的请求参数
type CreateRequest struct {
	// 分类名称
	Name string `json:"name" binding:"required,max=100" example:"工作"`
	// 分类描述
	Description string `json:"description" binding:"max=500" example:"工作相关的待办事项"`
	// 分类颜色
	Color string `json:"color" binding:"omitempty,len=7" example:"#FF0000"`
}

// CreateResponse 创建分类响应
// @Description 创建分类的响应
type CreateResponse struct {
	// 分类ID
	ID uint `json:"id" example:"1"`
}
