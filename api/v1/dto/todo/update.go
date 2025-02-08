package todo

// UpdateRequest 更新待办事项请求
// @Description 更新待办事项的请求参数
type UpdateRequest struct {
	// 标题
	Title string `json:"title,omitempty" binding:"omitempty,max=128" example:"完成项目文档"`
	// 描述
	Description string `json:"description,omitempty" binding:"omitempty,max=1024" example:"编写详细的项目设计文档"`
	// 状态
	Status string `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed" example:"in_progress"`
	// 优先级
	Priority string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high" example:"medium"`
	// 截止时间
	DueDate string `json:"dueDate,omitempty" example:"2024-02-08T17:12:40+08:00"`
	// 分类ID
	CategoryID *uint `json:"categoryId,omitempty" example:"1"`
	// 完成状态
	Completed *bool `json:"completed,omitempty" example:"false"`
}

// UpdateResponse 更新待办事项响应
type UpdateResponse struct {
	Message string `json:"message"` // 响应消息
}
