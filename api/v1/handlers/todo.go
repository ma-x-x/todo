package handlers

import (
	"net/http"
	"strconv"
	"todo-demo/internal/models"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateTodoRequest 定义创建待办事项的请求参数结构
type CreateTodoRequest struct {
	// Title 待办事项标题，必填，最大长度128个字符
	Title string `json:"title" binding:"required,max=128"`
	// Description 待办事项描述，可选，最大长度1024个字符
	Description string `json:"description" binding:"max=1024"`
	// Priority 优先级：1(低)/2(中)/3(高)
	Priority int `json:"priority" binding:"omitempty,oneof=1 2 3"`
	// CategoryID 所属分类ID，可选
	CategoryID *uint `json:"category_id" binding:"omitempty"`
}

// CreateTodoResponse 创建待办事项响应
// @Description 创建待办事项响应
type CreateTodoResponse struct {
	Todo *models.Todo `json:"todo"`
}

// UpdateTodoRequest 更新待办事项请求
// @Description 更新待办事项请求参数
type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"omitempty,max=128"`
	Description string `json:"description" binding:"omitempty,max=1024"`
	Completed   *bool  `json:"completed" binding:"omitempty"`
	Priority    int    `json:"priority" binding:"omitempty,oneof=1 2 3"`
	CategoryID  *uint  `json:"category_id"`
}

// UpdateTodoResponse 更新待办事项响应
// @Description 更新待办事项响应
type UpdateTodoResponse struct {
	Message string `json:"message"`
}

// ListTodoResponse 待办事项列表响应
// @Description 待办事项列表响应
type ListTodoResponse struct {
	Total int64          `json:"total"`
	Items []*models.Todo `json:"items"`
}

// CreateTodo 创建待办事项处理器
// @Summary 创建待办事项
// @Description 创建一个新的待办事项，可设置标题、描述、优先级和所属分类
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Param request body CreateTodoRequest true "创建待办事项的详细参数"
// @Success 200 {object} CreateTodoResponse "创建成功的待办事项信息"
// @Failure 400 {object} errors.Error "参数验证失败或业务错误"
// @Failure 401 {object} errors.Error "未授权访问"
// @Router /todos [post]
func CreateTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析并验证创建待办事项的请求参数
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 从上下文获取当前用户ID
		userID := c.GetUint("userID")
		// 调用业务层创建待办事项
		todo, err := todoService.Create(c.Request.Context(), userID, req.Title, req.Description,
			models.Priority(req.Priority), req.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 返回创建成功的待办事项信息
		c.JSON(http.StatusOK, CreateTodoResponse{
			Todo: todo,
		})
	}
}

// @Summary 获取待办事项列表
// @Description 获取当前用户的所有待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} ListTodoResponse
// @Failure 401 {object} errors.Error
// @Router /todos [get]
func ListTodos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		todos, err := todoService.List(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, ListTodoResponse{
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
// @Success 200 {object} models.Todo
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
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
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
// @Param request body UpdateTodoRequest true "更新待办事项请求参数"
// @Success 200 {object} UpdateTodoResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /todos/{id} [put]
func UpdateTodo(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		var req UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		userID := c.GetUint("userID")
		if err := todoService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, UpdateTodoResponse{
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
// @Success 200 {object} UpdateTodoResponse
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

		c.JSON(http.StatusOK, UpdateTodoResponse{
			Message: "Todo deleted successfully",
		})
	}
}
