package todo

// CreateRequest 创建待办事项的请求参数
type CreateRequest struct {
	// 待办事项标题
	// 必填，最大长度128个字符
	Title string `json:"title" binding:"required,max=128"`

	// 待办事项描述
	// 选填，最大长度1024个字符
	Description string `json:"description" binding:"max=1024"`

	// 分类ID
	// 选填，可以为空，表示未分类
	CategoryID *uint `json:"category_id"`

	// 优先级
	// 必填，可选值：
	// 0: 低优先级
	// 1: 中优先级
	// 2: 高优先级
	Priority int `json:"priority" binding:"oneof=0 1 2"`
}

// CreateResponse 创建待办事项的响应
type CreateResponse struct {
	// 新创建的待办事项ID
	ID uint `json:"id"`
}
