package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	Addr     string
	Password string
	Database int
}

func NewRedis(addr string,
	password string,
	database int) *RedisConnection {
	return &RedisConnection{
		Addr:     addr,
		Password: password,
		Database: database,
	}
}

func (r *RedisConnection) Connect(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password, // no password set
		DB:       r.Database, // use default DB
	})

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis connection was refused")
		return nil, err
	}
	fmt.Println(status)

	return rdb, nil
}
