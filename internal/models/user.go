package models

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
// 存储用户的基本信息，包括用户名、密码和邮箱
// 密码以加密形式存储，使用bcrypt加密算法
type User struct {
	Base
	Username string `json:"username" gorm:"uniqueIndex;size:32;not null"` // 用户名
	Password string `json:"-" gorm:"size:128;not null"`                   // 密码（加密存储）
	Email    string `json:"email" gorm:"size:128;uniqueIndex;not null"`   // 邮箱

	Categories []Category `json:"categories,omitempty" gorm:"foreignKey:UserID"` // 关联的分类列表
	Todos      []Todo     `json:"todos,omitempty" gorm:"foreignKey:UserID"`      // 关联的待办事项列表
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Validate 验证用户数据
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}

	if len(u.Username) < 3 || len(u.Username) > 32 {
		return errors.New("用户名长度应在3-32个字符之间")
	}

	if u.Email == "" {
		return errors.New("邮箱不能为空")
	}

	// 验证邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("邮箱格式无效")
	}

	if u.Password == "" {
		return errors.New("密码不能为空")
	}

	if len(u.Password) < 6 {
		return errors.New("密码长度不能小于6个字符")
	}

	return nil
}

// SetPassword 设置用户密码
//
// Parameters:
//   - password: 明文密码
//
// Returns:
//   - error: 返回加密过程中的错误，如果有的话
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度不能小于6个字符")
	}

	// 使用bcrypt算法对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证用户密码
//
// Parameters:
//   - password: 待验证的明文密码
//
// Returns:
//   - bool: 返回密码是否正确
func (u *User) CheckPassword(password string) bool {
	// 比较密码哈希值是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate 创建前的钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.Validate(); err != nil {
		return err
	}
	return u.Base.BeforeCreate(tx)
}

// BeforeUpdate 更新前的钩子函数
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if err := u.Validate(); err != nil {
		return err
	}
	return u.Base.BeforeUpdate(tx)
}
