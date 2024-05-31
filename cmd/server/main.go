package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jeanSagaz/rate-limiter/configs"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/internal/application/infra/web"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis/repository"
)

func main() {
	fmt.Println("rate-limiter")

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisConn := redis.NewRedis(configs.Addr, configs.Password, configs.Database)

	tokenConfiguration := []dto.TokenConfiguration{}
	err = json.Unmarshal([]byte(configs.TokenConfiguration), &tokenConfiguration)
	if err != nil {
		fmt.Println("json invalid format")
		panic(err)
	}

	ctx := context.Background()
	rdb, err := redisConn.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	redisRepository := repository.NewRedisRepositoryDb(rdb)
	h := web.NewHandler(redisRepository, tokenConfiguration, configs.NumberRequests, configs.Seconds)
	h.HandlerRequests()
}
