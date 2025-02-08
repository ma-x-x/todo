package handlers

import (
	"net/http"
	"strconv"
	"todo/api/v1/dto/category"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryResponse 分类响应
// @Description 分类信息响应
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

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器实例
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create 创建分类
// @Summary 创建分类
// @Description 创建一个新的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param request body category.CreateRequest true "创建分类请求参数"
// @Success 200 {object} response.Response{data=gin.H{id=int}} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req category.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	userID := middleware.GetUserID(c)
	id, err := h.categoryService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建分类失败")
		return
	}

	response.Success(c, gin.H{"id": id})
}

// Get 获取分类详情
// @Summary 获取分类详情
// @Description 获取指定ID的分类详情
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=CategoryResponse} "获取成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	userID := middleware.GetUserID(c)
	category, err := h.categoryService.Get(c.Request.Context(), uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类失败")
		return
	}

	// 转换为响应格式
	categoryResponse := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, categoryResponse)
}

// List 获取分类列表
// @Summary 获取分类列表
// @Description 获取当前用户的所有分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Success 200 {object} response.Response{data=gin.H{items=[]CategoryResponse,total=int}} "获取成功"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	categories, err := h.categoryService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	// 转换为响应格式
	categoryResponses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Color:       category.Color,
			CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response.Success(c, gin.H{
		"total": len(categories),
		"items": categoryResponses,
	})
}

// Update 更新分类
// @Summary 更新分类
// @Description 更新指定ID的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "分类ID"
// @Param request body category.UpdateRequest true "更新分类请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	var req category.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.categoryService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新分类失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除分类
// @Summary 删除分类
// @Description 删除指定ID的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.categoryService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除分类失败")
		return
	}

	response.Success(c, nil)
}
