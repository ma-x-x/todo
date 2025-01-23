package models

import (
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
// 存储用户的基本信息，包括用户名、密码和邮箱
// 密码以加密形式存储，使用bcrypt加密算法
type User struct {
	Base
	Username string `json:"username" gorm:"uniqueIndex;size:32"`
	Password string `json:"-" gorm:"size:128"`
	Email    string `json:"email" gorm:"size:128"`
}

// SetPassword 设置用户密码
// @param password 明文密码
// @return error 返回加密过程中的错误，如果有的话
func (u *User) SetPassword(password string) error {
	// 使用bcrypt算法对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证用户密码
// @param password 待验证的明文密码
// @return bool 返回密码是否正确
func (u *User) CheckPassword(password string) bool {
	// 比较密码哈希值是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
