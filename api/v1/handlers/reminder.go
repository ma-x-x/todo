package handlers

import (
	"net/http"
	"strconv"
	"todo-demo/api/v1/dto/reminder"
	"todo-demo/internal/models"
	"todo-demo/internal/service"
	"todo-demo/pkg/errors"

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
// @Success 200 {object} reminder.CreateResponse "创建成功的提醒信息"
// @Failure 400 {object} errors.Error "参数验证失败或业务错误"
// @Failure 401 {object} errors.Error "未授权访问"
// @Router /reminders [post]
func CreateReminder(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reminder.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 验证请求参数
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 转换提醒类型和通知类型
		reminderType, err := models.ParseReminderType(req.RemindType)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		notifyType, err := models.ParseNotifyType(req.NotifyType)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		userID := c.GetUint("userID")
		reminderModel, err := reminderService.Create(c.Request.Context(), userID, req.TodoID, req.RemindAt,
			reminderType, notifyType)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, reminder.CreateResponse{
			Reminder: reminderModel,
		})
	}
}

// @Summary 获取待办事项的提醒列表
// @Description 获取指定待办事项的所有提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param todo_id path int true "待办事项ID"
// @Success 200 {object} reminder.ListResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /reminders/todo/{todo_id} [get]
func ListReminders(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoID, err := strconv.ParseUint(c.Param("todo_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid todo ID"))
			return
		}

		reminders, err := reminderService.List(c.Request.Context(), uint(todoID))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, reminder.ListResponse{
			Total:     int64(len(reminders)),
			Reminders: reminders,
		})
	}
}

// @Summary 更新提醒
// @Description 更新指定的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "提醒ID"
// @Param request body reminder.UpdateRequest true "更新提醒请求参数"
// @Success 200 {object} reminder.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /reminders/{id} [put]
func UpdateReminder(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		var req reminder.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 验证请求参数
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		// 转换提醒类型和通知类型
		reminderType, err := models.ParseReminderType(req.RemindType)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		notifyType, err := models.ParseNotifyType(req.NotifyType)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		if err := reminderService.Update(c.Request.Context(), uint(id), req.RemindAt,
			reminderType, notifyType); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, reminder.UpdateResponse{
			Message: "Reminder updated successfully",
		})
	}
}

// @Summary 删除提醒
// @Description 删除指定的提醒
// @Tags 提醒管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int true "提醒ID"
// @Success 200 {object} reminder.UpdateResponse
// @Failure 400 {object} errors.Error
// @Failure 401 {object} errors.Error
// @Router /reminders/{id} [delete]
func DeleteReminder(reminderService service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, "Invalid ID"))
			return
		}

		if err := reminderService.Delete(c.Request.Context(), uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewError(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, reminder.UpdateResponse{
			Message: "Reminder deleted successfully",
		})
	}
}
