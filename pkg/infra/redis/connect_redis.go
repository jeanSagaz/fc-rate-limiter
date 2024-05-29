package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

func (d *RedisConnection) Set(c *redis.Client, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Set - error: ", err)
		return err
	}

	json := string(p)
	return c.Set(d.Context, key, json, duration).Err()
}

func (d *RedisConnection) Get(c *redis.Client, key string) (string, error) {
	value, err := c.Get(d.Context, key).Result()
	if err == redis.Nil {
		//if err == nil {
		fmt.Println("Get - error: ", err)
		return "", err
	}

	//var it := interface{}
	//json.Unmarshal([]byte(value), any)
	return value, nil
}
