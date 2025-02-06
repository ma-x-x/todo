// Package repository 实现数据访问层
package repository

import (
	"context"
	"todo/internal/models"
	"todo/pkg/errors"

	"gorm.io/gorm"
)

// UserRepository 定义用户仓储接口
type UserRepository interface {
	// Create 创建新用户
	// ctx: 上下文信息
	// user: 用户信息
	// 返回: error 创建过程中的错误信息
	Create(ctx context.Context, user *models.User) error

	// GetByID 根据用户ID获取用户信息
	// ctx: 上下文信息
	// id: 用户ID
	// 返回: (*models.User, error) 用户信息和可能的错误
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// GetByUsername 根据用户名获取用户信息
	// ctx: 上下文信息
	// username: 用户名
	// 返回: (*models.User, error) 用户信息和可能的错误
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// Update 更新用户信息
	// ctx: 上下文信息
	// user: 需要更新的用户信息
	// 返回: error 更新过程中的错误信息
	Update(ctx context.Context, user *models.User) error

	// Delete 删除用户
	// ctx: 上下文信息
	// id: 要删除的用户ID
	// 返回: error 删除过程中的错误信息
	Delete(ctx context.Context, id uint) error
}

// userRepo 实现 UserRepository 接口
type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
