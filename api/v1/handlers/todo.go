package handlers

import (
	"net/http"
	"strconv"
	"todo-demo/api/v1/dto/todo"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateTodo 创建待办事项处理器
// @Summary 创建待办事项
// @Description 创建一个新的待办事项，可设置标题、描述、优先级和所属分类
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Param request body todo.CreateRequest true "创建待办事项的详细参数，包含标题(必填)、描述(选填)、优先级(1低/2中/3高)和分类ID(选填)"
// @Success 200 {object} todo.CreateResponse "创建成功返回待办事项ID"
// @Failure 400 {object} errors.Error "参数验证失败或业务错误"
// @Failure 401 {object} errors.Error "未授权访问"
// @Router /todos [post]
func CreateTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req todo.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		userID := c.GetUint("userID")
		id, err := todoService.Create(c.Request.Context(), userID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, todo.CreateResponse{ID: id})
	}
}

// @Summary 获取待办事项列表
// @Description 获取当前用户的所有待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} todo.ListResponse
// @Failure 401 {object} errors.Error
// @Router /todos [get]
func ListTodos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		todos, err := todoService.List(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, todo.ListResponse{
			Total: int64(len(todos)),
			Items: todos,
		})
	}
}

// @Summary 获取待办事项详情
// @Description 获取指定的待办事项详情
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Success 200 {object} todo.DetailResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /todos/{id} [get]
func GetTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		todo, err := todoService.Get(c.Request.Context(), uint(id), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

// @Summary 更新待办事项
// @Description 更新指定的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Param request body todo.UpdateRequest true "更新待办事项请求参数"
// @Success 200 {object} todo.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /todos/{id} [put]
func UpdateTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req todo.UpdateRequest
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
		if err := todoService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, todo.UpdateResponse{
			Message: "Todo updated successfully",
		})
	}
}

// @Summary 删除待办事项
// @Description 删除指定的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "待办事项ID"
// @Success 200 {object} todo.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /todos/{id} [delete]
func DeleteTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := todoService.Delete(c.Request.Context(), uint(id), userID); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, todo.UpdateResponse{
			Message: "Todo deleted successfully",
		})
	}
}
