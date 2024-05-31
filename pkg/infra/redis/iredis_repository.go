package redis

import (
	"context"
	"time"
)

type IRedisRepository interface {
	SetValue(ctx context.Context, key string, value interface{}, duration time.Duration) error
	GetValue(ctx context.Context, key string) (string, error)
}
