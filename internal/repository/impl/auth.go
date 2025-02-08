package impl

import (
	"context"
	"fmt"
	"time"
	"todo/internal/repository/interfaces"
	"todo/pkg/errors"
	"todo/pkg/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuthRepository 认证仓储实现
type AuthRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewAuthRepository 创建认证仓储实例
func NewAuthRepository(db *gorm.DB, rdb *redis.Client) interfaces.AuthRepository {
	return &AuthRepository{
		db:  db,
		rdb: rdb,
	}
}

// SaveToken 保存用户令牌
func (r *AuthRepository) SaveToken(ctx context.Context, userID uint, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_token:%d", userID)
	return r.rdb.Set(ctx, key, token, expiration).Err()
}

// GetToken 获取用户令牌
func (r *AuthRepository) GetToken(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user_token:%d", userID)
	token, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.ErrTokenNotFound
	}
	return token, err
}

// DeleteToken 删除用户令牌
func (r *AuthRepository) DeleteToken(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user_token:%d", userID)
	return r.rdb.Del(ctx, key).Err()
}

// ValidateToken 验证令牌并返回用户ID
func (r *AuthRepository) ValidateToken(ctx context.Context, token string) (uint, error) {
	// 从token中解析用户ID
	userID, err := utils.ParseUserIDFromToken(token)
	if err != nil {
		return 0, errors.ErrInvalidToken
	}

	// 获取存储的token
	storedToken, err := r.GetToken(ctx, userID)
	if err != nil {
		return 0, err
	}

	// 验证token是否匹配
	if storedToken != token {
		return 0, errors.ErrInvalidToken
	}

	return userID, nil
}
