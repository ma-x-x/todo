package handlers

import (
	"log"
	"net/http"
	"strconv"
	"todo/api/v1/dto/category"
	"todo/api/v1/dto/reminder"
	"todo/api/v1/dto/todo"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// TodoHandler 待办事项处理器
type TodoHandler struct {
	todoService     service.TodoService
	categoryService service.CategoryService
}

// NewTodoHandler 创建待办事项处理器实例
func NewTodoHandler(todoService service.TodoService, categoryService service.CategoryService) *TodoHandler {
	return &TodoHandler{
		todoService:     todoService,
		categoryService: categoryService,
	}
}

// Create 创建待办事项
// @Summary 创建待办事项
// @Description 创建一个新的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param request body todo.CreateRequest true "创建待办事项请求参数"
// @Success 200 {object} response.Response{data=todo.TodoResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 403 {object} response.Response "分类不属于当前用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var req todo.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 如果指定了分类，验证分类是否存在且属于该用户
	if req.CategoryID != nil {
		userID := middleware.GetUserID(c)
		category, err := h.categoryService.Get(c.Request.Context(), *req.CategoryID, userID)
		if err != nil {
			log.Printf("验证分类失败: %v", err)
			response.BadRequest(c, "无效的分类ID")
			return
		}
		if category.UserID != userID {
			log.Printf("分类不属于当前用户: categoryID=%d, userID=%d", *req.CategoryID, userID)
			response.Error(c, http.StatusForbidden, "分类不属于当前用户")
			return
		}
	}

	userID := middleware.GetUserID(c)
	id, err := h.todoService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		log.Printf("创建待办事项失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "创建待办事项失败")
		return
	}

	log.Printf("成功创建待办事项: id=%d", id)

	// 获取创建后的完整信息
	createdTodo, err := h.todoService.Get(c.Request.Context(), id, userID)
	if err != nil {
		log.Printf("获取创建的待办事项失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "获取创建的待办事项失败")
		return
	}

	// 转换为响应格式
	todoResponse := todo.TodoResponse{
		ID:          createdTodo.ID,
		Title:       createdTodo.Title,
		Description: createdTodo.Description,
		Status:      createdTodo.Status.String(),
		Priority:    string(createdTodo.Priority),
		CategoryID:  createdTodo.CategoryID,
		Completed:   createdTodo.Completed,
		CreatedAt:   createdTodo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdTodo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Reminders:   make([]reminder.ReminderResponse, 0),
	}

	// 安全处理 DueDate
	if createdTodo.DueDate != nil {
		formattedDueDate := createdTodo.DueDate.Format("2006-01-02T15:04:05Z07:00")
		todoResponse.DueDate = formattedDueDate
	}

	if createdTodo.Category != nil {
		todoResponse.Category = &category.CategoryResponse{
			ID:          createdTodo.Category.ID,
			Name:        createdTodo.Category.Name,
			Description: createdTodo.Category.Description,
			Color:       createdTodo.Category.Color,
			CreatedAt:   createdTodo.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   createdTodo.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response.Success(c, todoResponse)
}

// List 获取待办事项列表
// @Summary 获取待办事项列表
// @Description 获取当前用户的所有待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Success 200 {object} response.Response{data=gin.H{items=[]todo.TodoResponse,total=int}} "获取成功"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /todos [get]
func (h *TodoHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	todos, err := h.todoService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待办事项列表失败")
		return
	}

	// 转换为响应格式
	todoResponses := make([]todo.TodoResponse, len(todos))
	for i, t := range todos {
		// 构建基本响应
		todoResponses[i] = todo.TodoResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status.String(),
			Priority:    string(t.Priority),
			CategoryID:  t.CategoryID,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Reminders:   make([]reminder.ReminderResponse, 0),
		}

		// 安全处理 DueDate
		if t.DueDate != nil {
			todoResponses[i].DueDate = t.DueDate.Format("2006-01-02T15:04:05Z07:00")
		}

		// 处理分类信息
		if t.Category != nil {
			todoResponses[i].Category = &category.CategoryResponse{
				ID:          t.Category.ID,
				Name:        t.Category.Name,
				Description: t.Category.Description,
				Color:       t.Category.Color,
				CreatedAt:   t.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:   t.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}

		// 处理提醒列表
		if len(t.Reminders) > 0 {
			reminderResps := make([]reminder.ReminderResponse, len(t.Reminders))
			for j, r := range t.Reminders {
				reminderResps[j] = reminder.ReminderResponse{
					ID:         r.ID,
					TodoID:     r.TodoID,
					RemindAt:   r.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
					RemindType: r.RemindType.String(),
					NotifyType: r.NotifyType.String(),
					Status:     r.Status,
					CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt:  r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				}
			}
			todoResponses[i].Reminders = reminderResps
		}
	}

	response.Success(c, gin.H{
		"total": len(todos),
		"items": todoResponses,
	})
}

// Get 获取待办事项详情
// @Summary 获取待办事项详情
// @Description 获取指定ID的待办事项详情
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=todo.TodoResponse} "获取成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /todos/{id} [get]
func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的待办事项ID")
		return
	}

	userID := middleware.GetUserID(c)
	t, err := h.todoService.Get(c.Request.Context(), uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待办事项失败")
		return
	}

	// 转换为响应格式
	todoResponse := todo.TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status.String(),
		Priority:    string(t.Priority),
		CategoryID:  t.CategoryID,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Reminders:   make([]reminder.ReminderResponse, 0),
	}

	// 安全处理 DueDate
	if t.DueDate != nil {
		todoResponse.DueDate = t.DueDate.Format("2006-01-02T15:04:05Z07:00")
	}

	// 处理分类信息
	if t.Category != nil {
		todoResponse.Category = &category.CategoryResponse{
			ID:          t.Category.ID,
			Name:        t.Category.Name,
			Description: t.Category.Description,
			Color:       t.Category.Color,
			CreatedAt:   t.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   t.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// 处理提醒列表
	if len(t.Reminders) > 0 {
		reminderResps := make([]reminder.ReminderResponse, len(t.Reminders))
		for i, r := range t.Reminders {
			reminderResps[i] = reminder.ReminderResponse{
				ID:         r.ID,
				TodoID:     r.TodoID,
				RemindAt:   r.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
				RemindType: r.RemindType.String(),
				NotifyType: r.NotifyType.String(),
				Status:     r.Status,
				CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:  r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		todoResponse.Reminders = reminderResps
	}

	response.Success(c, todoResponse)
}

// Update 更新待办事项
// @Summary 更新待办事项
// @Description 更新指定ID的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "待办事项ID"
// @Param request body todo.UpdateRequest true "更新待办事项请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的待办事项ID")
		return
	}

	var req todo.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.todoService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新待办事项失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除待办事项
// @Summary 删除待办事项
// @Description 删除指定ID的待办事项
// @Tags 待办事项管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的待办事项ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.todoService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除待办事项失败")
		return
	}

	response.Success(c, nil)
}
