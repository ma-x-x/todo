package errors

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// 标准错误定义
var (
	// 认证相关错误
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserExists         = errors.New("用户已存在")
	ErrInvalidToken       = errors.New("无效的令牌")
	ErrTokenExpired       = errors.New("令牌已过期")
	ErrTokenNotFound      = errors.New("令牌不存在")

	// Todo 相关错误
	ErrTodoNotFound     = errors.New("待办事项不存在")
	ErrCategoryNotFound = errors.New("分类不存在")
	ErrReminderNotFound = errors.New("提醒不存在")

	// 权限相关错误
	ErrUnauthorized = errors.New("未经授权的访问")
	ErrForbidden    = errors.New("禁止访问")
	ErrNoPermission = errors.New("无权限访问")

	// 数据验证错误
	ErrInvalidInput     = errors.New("无效的输入")
	ErrInvalidParameter = errors.New("无效的参数")

	// 数据库相关错误
	ErrDBConnection = errors.New("数据库连接失败")
	ErrDBQuery      = errors.New("数据库查询失败")
	ErrNotFound     = errors.New("record not found")
)

// APIError API错误响应结构体
// @Description API错误响应
type APIError struct {
	Code    int    `json:"code" example:"400"`                          // HTTP状态码
	Message string `json:"message" example:"Invalid request parameter"` // 错误信息
	Detail  string `json:"detail,omitempty"`                            // 详细错误信息（可选）
}

// Error 自定义错误类型
type Error struct {
	Message string
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}

// New 创建新的错误
func New(message string) error {
	return &Error{Message: message}
}

// Newf 创建带格式化的错误
func Newf(format string, args ...interface{}) error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

// NewAPIError 创建新的API错误响应
func NewAPIError(code int, message string, detail ...string) *APIError {
	err := &APIError{
		Code:    code,
		Message: message,
	}
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	return err
}

// ErrorCode 定义错误码类型
type ErrorCode int

const (
	// 系统级错误码
	ErrSystemCode ErrorCode = iota + 10000
	ErrDatabaseCode
	ErrCacheCode
	ErrNotFoundCode

	// 业务级错误码
	ErrInvalidAuthCode ErrorCode = iota + 20000
	ErrUserNotFoundCode
	ErrTodoNotFoundCode
)

// AppError 应用错误类型
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
	Stack   string
}

// Error 实现错误接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现错误解包
func (e *AppError) Unwrap() error {
	return e.Err
}

// WrapError 包装错误
func WrapError(err error, message string) *AppError {
	return &AppError{
		Message: message,
		Err:     err,
		Stack:   string(debug.Stack()),
	}
}

// WithCode 设置错误码
func (e *AppError) WithCode(code ErrorCode) *AppError {
	e.Code = code
	return e
}

// IsBusinessError 判断是否为业务错误
func (e *AppError) IsBusinessError() bool {
	return e.Code >= 20000 && e.Code < 30000
}

// IsSystemError 判断是否为系统错误
func (e *AppError) IsSystemError() bool {
	return e.Code >= 10000 && e.Code < 20000
}
