package todo

// UpdateRequest 更新待办事项请求
type UpdateRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,max=128"`       // 标题
	Description *string `json:"description,omitempty" binding:"omitempty,max=1024"` // 描述
	Completed   *bool   `json:"completed,omitempty"`                               // 完成状态
	Priority    *int    `json:"priority,omitempty" binding:"omitempty,oneof=1 2 3"` // 优先级
	CategoryID  *uint   `json:"categoryId,omitempty"`                              // 分类ID
}

// UpdateResponse 更新待办事项响应
type UpdateResponse struct {
	Message string `json:"message"` // 响应消息
}
