package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ParseValidationError 解析验证错误，返回友好的错误信息
func ParseValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Field() {
			case "Username":
				switch e.Tag() {
				case "required":
					return "请输入用户名"
				case "min":
					return "用户名长度不能少于3个字符"
				case "max":
					return "用户名长度不能超过32个字符"
				}
			case "Password":
				switch e.Tag() {
				case "required":
					return "请输入密码"
				case "min":
					return "密码长度不能少于6个字符"
				case "max":
					return "密码长度不能超过32个字符"
				}
			case "Email":
				switch e.Tag() {
				case "required":
					return "请输入邮箱"
				case "email":
					return "请输入有效的邮箱地址"
				}
			}
			// 如果没有匹配到具体的错误，返回通用错误信息
			return fmt.Sprintf("%s字段验证失败: %s", e.Field(), e.Tag())
		}
	}
	return ""
}
