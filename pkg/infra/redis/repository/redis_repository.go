package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	RedisDb *redis.Client
}

func NewRedisRepositoryDb(db *redis.Client) *RedisRepository {
	return &RedisRepository{RedisDb: db}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	json := string(b)
	return r.RedisDb.Set(ctx, key, json, duration).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.RedisDb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	}

	//var it := interface{}
	//json.Unmarshal([]byte(value), any)
	return value, nil
}
