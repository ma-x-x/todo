// Package response provides standard API response structures
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response format
// @Description API标准响应格式
type Response struct {
	// 状态码
	Code int `json:"code" example:"0"`

	// 响应消息
	Message string `json:"message" example:"success"`

	// 响应数据
	Data interface{} `json:"data,omitempty"`
}

// Success 创建成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 创建错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}
