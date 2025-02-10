package handlers

import (
	"net/http"
	"strconv"
	"todo/api/v1/dto/reminder"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReminderResponse 提醒响应
// @Description 提醒信息响应
type ReminderResponse struct {
	// 提醒ID
	ID uint `json:"id" example:"1"`
	// 待办事项ID
	TodoID uint `json:"todoId" example:"1"`
	// 提醒时间
	RemindAt string `json:"remindAt" example:"2024-02-08T17:10:54+08:00"`
	// 提醒类型
	RemindType string `json:"remindType" example:"once"`
	// 通知类型
	NotifyType string `json:"notifyType" example:"email"`
	// 提醒状态
	Status bool `json:"status" example:"false"`
	// 创建时间
	CreatedAt string `json:"createdAt" example:"2024-02-08T17:10:54+08:00"`
	// 更新时间
	UpdatedAt string `json:"updatedAt" example:"2024-02-08T17:10:54+08:00"`
}

// ReminderHandler 提醒处理器
type ReminderHandler struct {
	reminderService service.ReminderService
	todoService     service.TodoService
}

// NewReminderHandler 创建提醒处理器实例
func NewReminderHandler(
	reminderService service.ReminderService,
	todoService service.TodoService,
) *ReminderHandler {
	return &ReminderHandler{
		reminderService: reminderService,
		todoService:     todoService,
	}
}

// Create 创建提醒
// @Summary 创建提醒
// @Description 为待办事项创建一个新的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param request body reminder.CreateRequest true "创建提醒请求参数"
// @Success 200 {object} response.Response{data=gin.H{id=int}} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /reminders [post]
func (h *ReminderHandler) Create(c *gin.Context) {
	var req reminder.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	id, err := h.reminderService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建提醒失败")
		return
	}

	response.Success(c, gin.H{"id": id})
}

// List 获取提醒列表
// @Summary 获取提醒列表
// @Description 获取指定待办事项的所有提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param todo_id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=gin.H{items=[]ReminderResponse,total=int}} "获取成功"
// @Failure 400 {object} response.Response "无效的待办事项ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /reminders/todo/{todo_id} [get]
func (h *ReminderHandler) List(c *gin.Context) {
	todoID, err := strconv.ParseUint(c.Param("todo_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的待办事项ID")
		return
	}

	// 获取用户ID
	userID := middleware.GetUserID(c)

	// 验证待办事项是否属于当前用户
	todo, err := h.todoService.Get(c.Request.Context(), uint(todoID), userID)
	if err != nil {
		response.Error(c, http.StatusForbidden, "待办事项不存在或无权访问")
		return
	}

	reminders, err := h.reminderService.ListByTodoID(c.Request.Context(), todo.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取提醒列表失败")
		return
	}

	// 转换为响应格式
	reminderResponses := make([]ReminderResponse, 0) // 初始化为空数组
	for _, reminder := range reminders {
		reminderResponses = append(reminderResponses, ReminderResponse{
			ID:         reminder.ID,
			TodoID:     reminder.TodoID,
			RemindAt:   reminder.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
			RemindType: reminder.RemindType.String(),
			NotifyType: reminder.NotifyType.String(),
			Status:     reminder.Status,
			CreatedAt:  reminder.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  reminder.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	response.Success(c, gin.H{
		"total": len(reminderResponses),
		"items": reminderResponses,
	})
}

// Update 更新提醒
// @Summary 更新提醒
// @Description 更新指定ID的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "提醒ID"
// @Param request body reminder.UpdateRequest true "更新提醒请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /reminders/{id} [put]
func (h *ReminderHandler) Update(c *gin.Context) {
	var req reminder.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的提醒ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.reminderService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新提醒失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除提醒
// @Summary 删除提醒
// @Description 删除指定ID的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "提醒ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的提醒ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /reminders/{id} [delete]
func (h *ReminderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的提醒ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.reminderService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除提醒失败")
		return
	}

	response.Success(c, nil)
}

// Get 获取提醒详情
// @Summary 获取提醒详情
// @Description 获取指定ID的提醒详情
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param id path int true "提醒ID"
// @Success 200 {object} response.Response{data=ReminderResponse} "获取成功"
// @Failure 400 {object} response.Response "无效的提醒ID"
// @Failure 401 {object} response.Response "未授权访问"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /reminders/{id} [get]
func (h *ReminderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的提醒ID")
		return
	}

	userID := middleware.GetUserID(c)
	reminder, err := h.reminderService.Get(c.Request.Context(), uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取提醒失败")
		return
	}

	reminderResponse := ReminderResponse{
		ID:         reminder.ID,
		TodoID:     reminder.TodoID,
		RemindAt:   reminder.RemindAt.Format("2006-01-02T15:04:05Z07:00"),
		RemindType: reminder.RemindType.String(),
		NotifyType: reminder.NotifyType.String(),
		Status:     reminder.Status,
		CreatedAt:  reminder.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  reminder.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, reminderResponse)
}
