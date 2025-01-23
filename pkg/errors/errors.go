package errors

import "errors"

var (
	// Auth errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")

	// Todo errors
	ErrTodoNotFound     = errors.New("todo not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrReminderNotFound = errors.New("reminder not found")

	// Permission errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Validation errors
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidParameter = errors.New("invalid parameter")
)

// Error 错误响应
// @Description API错误响应
type Error struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request parameter"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
