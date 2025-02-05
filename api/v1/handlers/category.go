package handlers

import (
	"net/http"
	"strconv"
	"todo-demo/api/v1/dto/category"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateCategory 创建分类处理器
// @Summary 创建分类
// @Description 创建一个新的待办事项分类，包含名称和可选的颜色属性
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Param request body category.CreateRequest true "创建分类的请求参数"
// @Success 200 {object} category.CreateResponse "创建成功返回的分类信息"
// @Failure 400 {object} errors.Error "请求参数错误"
// @Failure 401 {object} errors.Error "未授权访问"
// @Router /categories [post]
func CreateCategory(categoryService service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req category.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		userID := c.GetUint("userID")
		id, err := categoryService.Create(c.Request.Context(), userID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, category.CreateResponse{
			ID: id,
		})
	}
}

// @Summary 获取分类列表
// @Description 获取当前用户的所有分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} category.ListResponse
// @Failure 401 {object} errors.Error
// @Router /categories [get]
func ListCategories(categoryService service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		categories, err := categoryService.List(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, category.ListResponse{
			Total: int64(len(categories)),
			Items: categories,
		})
	}
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 更新指定的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "分类ID"
// @Param request body category.UpdateRequest true "更新分类请求参数"
// @Success 200 {object} category.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /categories/{id} [put]
func UpdateCategory(categoryService service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req category.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := categoryService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
	}
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 删除指定的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "分类ID"
// @Success 200 {object} category.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /categories/{id} [delete]
func DeleteCategory(categoryService service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := categoryService.Delete(c.Request.Context(), uint(id), userID); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, category.UpdateResponse{
			Message: "Category deleted successfully",
		})
	}
}
