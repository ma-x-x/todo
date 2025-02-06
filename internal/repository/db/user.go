// Package db 提供数据库访问的具体实现
package db

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// userRepository 实现用户数据库操作的结构体
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储的实例
// db: 数据库连接实例
// 返回: 用户仓储实例
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// Create 在数据库中创建新的用户记录
// ctx: 上下文信息
// user: 要创建的用户信息
// 返回: error 创建过程中的错误信息
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID从数据库获取用户信息
// ctx: 上下文信息
// id: 用户ID
// 返回: (*models.User, error) 用户信息和可能的错误
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名从数据库获取用户信息
// ctx: 上下文信息
// username: 用户名
// 返回: (*models.User, error) 用户信息和可能的错误
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新数据库中的用户记录
// ctx: 上下文信息
// user: 需要更新的用户信息
// 返回: error 更新过程中的错误信息
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 从数据库中删除用户记录
// ctx: 上下文信息
// id: 要删除的用户ID
// 返回: error 删除过程中的错误信息
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
