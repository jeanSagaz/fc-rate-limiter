package repository

import (
	"context"
	"time"
)

type IRedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
