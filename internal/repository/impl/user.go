package impl

import (
	"context"
	"todo/internal/models"
	"todo/internal/repository/interfaces"

	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Create(ctx, user)
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, r.handleError(err, "user")
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, r.handleError(err, "user")
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.GetDB(ctx).Where("email = ?", email).First(&user).Error
	return &user, r.handleError(err, "user")
}

// Update 更新用户信息
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Update(ctx, user)
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &models.User{Base: models.Base{ID: id}})
}

// List 获取用户列表
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := r.BaseRepository.List(ctx, offset, limit, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Count 获取用户总数
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.BaseRepository.Count(ctx, &models.User{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
