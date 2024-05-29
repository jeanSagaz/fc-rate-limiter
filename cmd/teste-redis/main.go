package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
)

type Person struct {
	// Name string `redis:"name"`
	// Age  int    `redis:"age"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var cache = redis.RedisConnection{
	Addr:     "localhost:6379",
	Password: "Redis",
	Database: 0,
}

func main() {
	fmt.Println("rate-limiter")

	var person = Person{"Morpheus", 35}
	ctx := context.TODO()
	cache.Context = ctx
	//redis.ConnectRedis(ctx)

	rdb, _ := cache.Connect(ctx)
	cache.Set(rdb, "TESTE", person, time.Minute*1)
	value, _ := cache.Get(rdb, "TESTE")
	fmt.Println("value: " + value)
}
