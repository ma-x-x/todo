// Package response provides standard API response structures
package response

import (
	"net/http"
)

// Response represents a standard API response format
// @Description 标准API响应结构
type Response struct {
	// HTTP状态码
	Code    int         `json:"code" example:"200"`
	
	// 响应消息
	Message string      `json:"message" example:"Success"`
	
	// 响应数据
	Data    interface{} `json:"data,omitempty"`
	
	// 请求追踪ID，用于调试
	TraceID string      `json:"traceId,omitempty"`
}

// NewResponse creates a new Response instance
// NewResponse 创建新的响应实例
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// WithTraceID adds a trace ID to the response
// WithTraceID 添加追踪ID到响应中
func (r *Response) WithTraceID(traceID string) *Response {
	r.TraceID = traceID
	return r
}

const (
	// 成功响应的默认消息
	SuccessMsg = "操作成功"
	// 默认错误消息
	DefaultErrorMsg = "服务器内部错误"
)

// Success creates a success response
// Success 创建成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: SuccessMsg,
		Data:    data,
	}
}

// Error creates an error response
// Error 创建错误响应
func Error(code int, message string) *Response {
	if message == "" {
		message = DefaultErrorMsg
	}
	return &Response{
		Code:    code,
		Message: message,
	}
} 