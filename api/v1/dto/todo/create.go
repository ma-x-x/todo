package todo

// CreateRequest 创建待办事项请求
// @Description 创建待办事项的请求参数
type CreateRequest struct {
	// Title 待办事项标题
	// Required: true
	// Max Length: 128
	Title string `json:"title" binding:"required,max=128" example:"完成项目文档"`

	// Description 待办事项描述
	// Required: false
	// Max Length: 1024
	Description string `json:"description" binding:"max=1024" example:"编写详细的项目设计文档"`

	// Priority 优先级
	// Required: false
	// Enum: [low medium high]
	// Example: medium
	Priority string `json:"priority" binding:"omitempty,oneof=low medium high" example:"medium"`

	// DueDate 截止时间
	// Required: false
	// Example: 2024-02-08T17:12:40+08:00
	DueDate string `json:"dueDate" example:"2024-02-08T17:12:40+08:00"`

	// CategoryID 所属分类ID
	// Required: false
	CategoryID *uint `json:"categoryId" binding:"omitempty" example:"1"`
}

// CreateResponse 创建待办事项响应
type CreateResponse struct {
	// ID 新创建的待办事项ID
	ID uint `json:"id"`
}
