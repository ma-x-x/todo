package handlers

import (
	"net/http"
	"strconv"
	"todo/api/v1/dto/reminder"
	"todo/internal/service"
	"todo/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateReminder 创建提醒处理器
// @Summary 创建提醒
// @Description 为待办事项创建定时提醒，支持单次、每日和每周重复的提醒方式
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT认证令牌"
// @Param request body reminder.CreateRequest true "创建提醒的详细参数"
// @Success 200 {object} response.Response{data=reminder.CreateResponse} "创建成功的提醒信息"
// @Failure 400 {object} response.Response "参数验证失败或业务错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /reminders [post]
func CreateReminder(reminderService service.ReminderService, todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reminder.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
			return
		}

		userID := c.GetUint("userID")
		id, err := reminderService.Create(c.Request.Context(), userID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		// 获取创建后的完整提醒信息
		createdReminder, err := reminderService.Get(c.Request.Context(), id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "获取创建的提醒失败"))
			return
		}

		// 获取关联的待办事项信息
		todo, err := todoService.Get(c.Request.Context(), req.TodoID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "获取关联的待办事项失败"))
			return
		}

		resp := reminder.CreateResponse{
			ID:         createdReminder.ID,
			TodoID:     createdReminder.TodoID,
			RemindAt:   createdReminder.RemindAt,
			RemindType: createdReminder.RemindType,
			NotifyType: createdReminder.NotifyType,
			CreatedAt:  createdReminder.CreatedAt,
			Todo:       todo,
			Reminder:   createdReminder,
		}

		c.JSON(http.StatusOK, response.Success(resp))
	}
}

// ListReminders 获取提醒列表
// @Summary 获取待办事项的提醒列表
// @Description 获取指定待办事项的所有提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param todo_id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=reminder.ListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /reminders/todo/{todo_id} [get]
func ListReminders(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoID, err := strconv.ParseUint(c.Param("todo_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid todo ID"))
			return
		}

		reminders, err := reminderService.ListByTodoID(c.Request.Context(), uint(todoID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(reminder.ListResponse{
			Items: reminders,
			Total: int64(len(reminders)),
		}))
	}
}

// UpdateReminder 更新提醒
// @Summary 更新提醒
// @Description 更新指定的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "提醒ID"
// @Param request body reminder.UpdateRequest true "更新提醒请求参数"
// @Success 200 {object} response.Response{data=reminder.UpdateResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /reminders/{id} [put]
func UpdateReminder(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reminder.UpdateRequest
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
		if err := reminderService.Update(c.Request.Context(), uint(id), userID, &req); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		// 获取更新后的完整数据
		updatedReminder, err := reminderService.Get(c.Request.Context(), uint(id), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch updated reminder"))
			return
		}

		c.JSON(http.StatusOK, response.Success(updatedReminder))
	}
}

// DeleteReminder 删除提醒
// @Summary 删除提醒
// @Description 删除指定的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "提醒ID"
// @Success 200 {object} response.Response{data=reminder.UpdateResponse} "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权访问"
// @Router /reminders/{id} [delete]
func DeleteReminder(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid ID"))
			return
		}

		userID := c.GetUint("userID")
		if err := reminderService.Delete(c.Request.Context(), uint(id), userID); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(reminder.UpdateResponse{
			Message: "Reminder deleted successfully",
		}))
	}
}
