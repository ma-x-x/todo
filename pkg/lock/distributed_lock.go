package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributedLock struct {
	rdb *redis.Client
}

func NewDistributedLock(rdb *redis.Client) *DistributedLock {
	return &DistributedLock{rdb: rdb}
}

func (l *DistributedLock) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return l.rdb.SetNX(ctx, "lock:"+key, "1", expiration).Result()
}

func (l *DistributedLock) Unlock(ctx context.Context, key string) error {
	return l.rdb.Del(ctx, "lock:"+key).Err()
}
