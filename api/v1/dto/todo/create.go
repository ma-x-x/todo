package todo

// CreateRequest 创建待办事项请求
type CreateRequest struct {
	// Title 待办事项标题
	// Required: true
	// Max Length: 128
	Title       string `json:"title" binding:"required,max=128"`

	// Description 待办事项描述
	// Required: false
	// Max Length: 1024
	Description string `json:"description" binding:"max=1024"`

	// Priority 优先级
	// Required: false
	// Enum: [1,2,3]
	// Example: 2
	// Note: 1=低优先级 2=中优先级 3=高优先级
	Priority    int    `json:"priority" binding:"omitempty,oneof=1 2 3"`

	// CategoryID 所属分类ID
	// Required: false
	CategoryID  *uint  `json:"categoryId" binding:"omitempty"`
}

// CreateResponse 创建待办事项响应
type CreateResponse struct {
	// ID 新创建的待办事项ID
	ID uint `json:"id"`
}
