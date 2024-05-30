package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisRepository interface {
	Set(c *redis.Client, key string, value interface{}, duration time.Duration) error
	Get(c *redis.Client, key string) (string, error)
}
