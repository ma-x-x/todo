package middleware

import (
	"time"
	"todo/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// MiddlewareChain 中间件链管理
type MiddlewareChain struct {
	middlewares []gin.HandlerFunc
	order       []int
}

// NewMiddlewareChain 创建新的中间件链
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]gin.HandlerFunc, 0),
		order:       make([]int, 0),
	}
}

// Add 添加中间件
func (mc *MiddlewareChain) Add(middleware gin.HandlerFunc, order int) {
	idx := 0
	for i, o := range mc.order {
		if order < o {
			break
		}
		idx = i + 1
	}

	mc.middlewares = append(mc.middlewares, nil)
	mc.order = append(mc.order, 0)
	copy(mc.middlewares[idx+1:], mc.middlewares[idx:])
	copy(mc.order[idx+1:], mc.order[idx:])
	mc.middlewares[idx] = middleware
	mc.order[idx] = order
}

// Apply 应用中间件链到路由组
func (mc *MiddlewareChain) Apply(group *gin.RouterGroup) {
	for _, middleware := range mc.middlewares {
		group.Use(middleware)
	}
}

// MiddlewareContext 中间件上下文
type MiddlewareContext struct {
	TraceID   string
	StartTime time.Time
	UserID    uint
	RequestID string
}

// NewContext 创建新的中间件上下文
func NewContext() *MiddlewareContext {
	return &MiddlewareContext{
		TraceID:   uuid.New().String(),
		StartTime: time.Now(),
		RequestID: uuid.New().String(),
	}
}

// GetContext 从gin上下文获取中间件上下文
func GetContext(c *gin.Context) *MiddlewareContext {
	if ctx, exists := c.Get("middleware_context"); exists {
		return ctx.(*MiddlewareContext)
	}
	return NewContext()
}

// SetContext 设置中间件上下文到gin上下文
func SetContext(c *gin.Context, ctx *MiddlewareContext) {
	c.Set("middleware_context", ctx)
}

// CommonMiddlewares 通用中间件集合
type CommonMiddlewares struct {
	middlewares []gin.HandlerFunc
}

// NewCommonMiddlewares 创建通用中间件集合
func NewCommonMiddlewares(cfg *config.Config, rdb *redis.Client) *CommonMiddlewares {
	return &CommonMiddlewares{
		middlewares: []gin.HandlerFunc{
			gin.Logger(),
			gin.Recovery(),
			Cors(),
			RateLimiterMiddleware(rdb, cfg.RateLimit.RequestsPerSecond, time.Second),
			MetricsMiddleware(),
			LoggerMiddleware(),
			RequestIDMiddleware(),
			TraceMiddleware(),
			CacheMiddleware(rdb, 5*time.Minute),
		},
	}
}

// Apply 应用中间件到路由组
func (m *CommonMiddlewares) Apply(group *gin.RouterGroup) {
	for _, middleware := range m.middlewares {
		group.Use(middleware)
	}
}

// AuthMiddlewares 认证中间件集合
type AuthMiddlewares struct {
	middlewares []gin.HandlerFunc
}

// NewAuthMiddlewares 创建认证中间件集合
func NewAuthMiddlewares(cfg *config.Config) *AuthMiddlewares {
	return &AuthMiddlewares{
		middlewares: []gin.HandlerFunc{
			Auth(&cfg.JWT),
		},
	}
}

// Apply 应用中间件到路由组
func (m *AuthMiddlewares) Apply(group *gin.RouterGroup) {
	for _, middleware := range m.middlewares {
		group.Use(middleware)
	}
}
