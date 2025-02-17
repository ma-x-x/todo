package todo

import (
	"todo/api/v1/dto/category"
	"todo/api/v1/dto/reminder"
)

// TodoResponse Todo响应
type TodoResponse struct {
	// Todo ID
	ID uint `json:"id" example:"1"`
	// 标题
	Title string `json:"title" example:"完成项目文档"`
	// 描述
	Description string `json:"description" example:"编写详细的项目设计文档"`
	// 状态
	Status string `json:"status" example:"pending"`
	// 优先级
	Priority string `json:"priority" example:"medium"`
	// 截止时间
	DueDate string `json:"dueDate,omitempty" example:"2024-02-08T17:12:40+08:00"`
	// 分类ID
	CategoryID *uint `json:"categoryId,omitempty" example:"1"`
	// 完成状态
	Completed bool `json:"completed" example:"false"`
	// 创建时间
	CreatedAt string `json:"createdAt" example:"2024-02-08T17:12:40+08:00"`
	// 更新时间
	UpdatedAt string `json:"updatedAt" example:"2024-02-08T17:12:40+08:00"`
	// 分类信息
	Category *category.CategoryResponse `json:"category,omitempty"`
	// 提醒列表
	Reminders []reminder.ReminderResponse `json:"reminders"`
}
