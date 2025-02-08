package interfaces

import (
	"context"
	"time"
)

// AuthRepository 认证仓储接口
type AuthRepository interface {
	// SaveToken 保存用户令牌
	SaveToken(ctx context.Context, userID uint, token string, expiration time.Duration) error

	// GetToken 获取用户令牌
	GetToken(ctx context.Context, userID uint) (string, error)

	// DeleteToken 删除用户令牌
	DeleteToken(ctx context.Context, userID uint) error

	// ValidateToken 验证令牌并返回用户ID
	ValidateToken(ctx context.Context, token string) (uint, error)
}
