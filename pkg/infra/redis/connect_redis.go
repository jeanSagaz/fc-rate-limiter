package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	Context  context.Context
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

func (d *RedisConnection) Connect(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     d.Addr,
		Password: d.Password, // no password set
		DB:       d.Database, // use default DB
	})

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Redis connection was refused")
		return nil, err
	}
	fmt.Println(status)

	return rdb, nil
}
