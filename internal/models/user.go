package models

import (
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
// @Description 用户信息
type User struct {
	Base
	Username string `json:"username" gorm:"size:32;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:128;not null"`
	Email    string `json:"email" gorm:"size:128;not null"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
