package todo

// UpdateRequest 更新待办事项的请求参数
type UpdateRequest struct {
	// 待办事项标题
	// 选填，最大长度128个字符
	// 如果不需要更新，可以不传或传null
	Title *string `json:"title" binding:"omitempty,max=128"`

	// 待办事项描述
	// 选填，最大长度1024个字符
	// 如果不需要更新，可以不传或传null
	Description *string `json:"description" binding:"omitempty,max=1024"`

	// 是否完成
	// 选填，true表示已完成，false表示未完成
	// 如果不需要更新，可以不传或传null
	Completed *bool `json:"completed"`

	// 分类ID
	// 选填，可以为空表示取消分类
	// 如果不需要更新，可以不传或传null
	CategoryID *uint `json:"category_id"`

	// 优先级
	// 选填，可选值：
	// 0: 低优先级
	// 1: 中优先级
	// 2: 高优先级
	// 如果不需要更新，可以不传或传null
	Priority *int `json:"priority" binding:"omitempty,oneof=0 1 2"`
}

// UpdateResponse 更新待办事项的响应
type UpdateResponse struct {
	// 更新结果消息
	Message string `json:"message"`
}
