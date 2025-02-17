package category

// CategoryResponse 分类响应
type CategoryResponse struct {
	// 分类ID
	ID uint `json:"id" example:"1"`
	// 分类名称
	Name string `json:"name" example:"工作"`
	// 分类描述
	Description string `json:"description" example:"工作相关的待办事项"`
	// 分类颜色
	Color string `json:"color" example:"#FF0000"`
	// 创建时间
	CreatedAt string `json:"createdAt" example:"2024-02-08T17:10:54+08:00"`
	// 更新时间
	UpdatedAt string `json:"updatedAt" example:"2024-02-08T17:10:54+08:00"`
}
