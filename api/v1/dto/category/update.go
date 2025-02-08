package category

import (
	"errors"
)

// UpdateRequest 更新分类请求
// @Description 更新分类的请求参数
type UpdateRequest struct {
	// 分类名称
	// Max length: 100
	Name string `json:"name,omitempty" binding:"omitempty,max=100" example:"工作"`
	// 分类描述
	// Max length: 500
	Description string `json:"description,omitempty" binding:"omitempty,max=500" example:"工作相关的待办事项"`
	// 分类颜色
	// Pattern: ^#[0-9A-Fa-f]{6}$
	Color string `json:"color,omitempty" binding:"omitempty,len=7,hexcolor" example:"#FF0000"`
}

// UpdateResponse 更新分类响应
// @Description 更新分类的响应
type UpdateResponse struct {
	// 响应消息
	Message string `json:"message" example:"更新成功"`
}

// Validate 验证请求参数
func (r *UpdateRequest) Validate() error {
	if r.Name == "" && r.Description == "" && r.Color == "" {
		return errors.New("至少需要更新一个字段")
	}
	return nil
}
