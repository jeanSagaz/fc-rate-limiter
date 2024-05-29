package main

import (
	"encoding/json"
	"fmt"

	"github.com/jeanSagaz/rate-limiter/configs"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/internal/web"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
)

func main() {
	fmt.Println("rate-limiter")

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisConfig := redis.NewRedis(configs.Addr, configs.Password, configs.Database)

	tokenConfiguration := []dto.TokenConfiguration{}
	err = json.Unmarshal([]byte(configs.TokenConfiguration), &tokenConfiguration)
	if err != nil {
		fmt.Println("json invalid format")
		panic(err)
	}

	h := web.NewHandler(redisConfig, tokenConfiguration, configs.NumberRequests, configs.Seconds)
	h.HandlerRequests()
}