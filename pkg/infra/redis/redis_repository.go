package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func (d *RedisConnection) Set(c *redis.Client, key string, value interface{}, duration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	json := string(b)
	return c.Set(d.Context, key, json, duration).Err()
}

func (d *RedisConnection) Get(c *redis.Client, key string) (string, error) {
	value, err := c.Get(d.Context, key).Result()
	if err == redis.Nil {
		return "", err
	}

	//var it := interface{}
	//json.Unmarshal([]byte(value), any)
	return value, nil
}
