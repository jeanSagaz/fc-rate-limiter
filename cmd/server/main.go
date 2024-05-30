package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jeanSagaz/rate-limiter/configs"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/internal/web"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
	"github.com/spf13/viper"
)

type Conf struct {
	Addr               string `mapstructure:"ADDR"`
	Password           string `mapstructure:"PASSWORD"`
	Database           int    `mapstructure:"DATABASE"`
	TokenConfiguration string `mapstructure:"TOKEN_CONFIGURATION"`
	NumberRequests     int    `mapstructure:"NUMBER_REQUESTS"`
	Seconds            int    `mapstructure:"SECONDS"`
}

func loadConfig(path string) (*Conf, error) {
	var cfg *Conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	// viper.AddConfigPath(path)
	fmt.Println(path)
	viper.SetConfigFile(".env")
	viper.SetConfigFile("./cmd/server/.env")
	//viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}

func main() {
	fmt.Println("rate-limiter")

	configs, err := configs.LoadConfig(".")
	//configs, err := loadConfig("./cmd/server/.env")
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

	redisRepository := redis.NewRedisRepositoryDb(rdb)
	h := web.NewHandler(redisRepository, tokenConfiguration, configs.NumberRequests, configs.Seconds)
	h.HandlerRequests()
}
