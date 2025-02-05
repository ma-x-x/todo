package errors

import "errors"

var (
	// 认证相关错误
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserExists         = errors.New("用户已存在")
	ErrInvalidToken       = errors.New("无效的令牌")
	ErrTokenExpired       = errors.New("令牌已过期")

	// Todo 相关错误
	ErrTodoNotFound     = errors.New("待办事项不存在")
	ErrCategoryNotFound = errors.New("分类不存在")
	ErrReminderNotFound = errors.New("提醒不存在")

	// 权限相关错误
	ErrUnauthorized = errors.New("未经授权的访问")
	ErrForbidden    = errors.New("禁止访问")

	// 数据验证错误
	ErrInvalidInput     = errors.New("无效的输入")
	ErrInvalidParameter = errors.New("无效的参数")

	// 数据库相关错误
	ErrDBConnection = errors.New("数据库连接失败")
	ErrDBQuery      = errors.New("数据库查询失败")
)

// Error 错误响应结构体
// @Description API错误响应
type Error struct {
	Code    int    `json:"code" example:"400"`                          // HTTP状态码
	Message string `json:"message" example:"Invalid request parameter"` // 错误信息
	Detail  string `json:"detail,omitempty"`                            // 详细错误信息（可选）
}

// NewError 创建新的错误响应
//
// Parameters:
//   - code: HTTP状态码
//   - message: 错误信息
//   - detail: 详细错误信息（可选）
//
// Returns:
//   - *Error: 返回错误响应对象
func NewError(code int, message string, detail ...string) *Error {
	err := &Error{
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
	
	// 业务级错误码
	ErrInvalidAuthCode ErrorCode = iota + 20000
	ErrUserNotFoundCode
	ErrTodoNotFoundCode
)

// 将错误码映射到错误信息
var errorMessages = map[ErrorCode]string{
	ErrSystemCode:       "系统错误",
	ErrDatabaseCode:     "数据库错误",
	ErrCacheCode:        "缓存错误",
	ErrInvalidAuthCode:  "认证失败",
	ErrUserNotFoundCode: "用户不存在",
	ErrTodoNotFoundCode: "待办事项不存在",
}
