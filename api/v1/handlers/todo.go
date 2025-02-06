package handlers

import (
	"net/http"
	"strconv"
	"todo/api/v1/dto/todo"
	"todo/internal/service"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateTodo 创建待办事项处理器
// @Summary 创建待办事项
// @Description 创建一个新的待办事项，可设置标题、描述、优先级和所属分类
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Param request body todo.CreateRequest true "创建待办事项的详细参数"
// @Success 200 {object} response.Response{data=todo.DetailResponse} "创建成功返回待办事项信息"
// @Failure 400 {object} response.Response "参数验证失败或业务错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /todos [post]
func CreateTodo(todoService service.TodoService, categoryService service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req todo.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
			return
		}

		// 如果提供了 CategoryID，先验证分类是否存在
		if req.CategoryID != nil {
			// 验证分类是否存在且属于当前用户
			category, err := categoryService.Get(c.Request.Context(), *req.CategoryID, c.GetUint("userID"))
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid category ID"))
				return
			}
			if category.UserID != c.GetUint("userID") {
				c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Category does not belong to user"))
				return
			}
		}

		userID := c.GetUint("userID")
		id, err := todoService.Create(c.Request.Context(), userID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		// 获取创建后的完整数据
		createdTodo, err := todoService.Get(c.Request.Context(), id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch created todo"))
			return
		}

		c.JSON(http.StatusOK, response.Success(createdTodo))
	}
}

// ListTodos 获取待办事项列表
// @Summary 获取待办事项列表
// @Description 获取当前用户的所有待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} response.Response{data=todo.ListResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /todos [get]
func ListTodos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		todos, err := todoService.List(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(todo.ListResponse{
			Total: int64(len(todos)),
			Items: todos,
		}))
	}
}

// GetTodo 获取待办事项详情
// @Summary 获取待办事项详情
// @Description 获取指定的待办事项详情
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=todo.DetailResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /todos/{id} [get]
func GetTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		todoItem, err := todoService.Get(c.Request.Context(), uint(id), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(todo.DetailResponse{
			Todo: todoItem,
		}))
	}
}

// UpdateTodo 更新待办事项
// @Summary 更新待办事项
// @Description 更新指定的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Param request body todo.UpdateRequest true "更新待办事项请求参数"
// @Success 200 {object} response.Response{data=todo.DetailResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /todos/{id} [put]
func UpdateTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req todo.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := todoService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		// 获取更新后的完整数据
		updatedTodo, err := todoService.Get(c.Request.Context(), uint(id), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch updated todo"))
			return
		}

		c.JSON(http.StatusOK, response.Success(updatedTodo))
	}
}

// DeleteTodo 删除待办事项
// @Summary 删除待办事项
// @Description 删除指定的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=todo.UpdateResponse} "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /todos/{id} [delete]
func DeleteTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := todoService.Delete(c.Request.Context(), uint(id), userID); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(todo.UpdateResponse{
			Message: "Todo deleted successfully",
		}))
	}
}
