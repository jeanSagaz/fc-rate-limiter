package web

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis/repository"
	"github.com/stretchr/testify/require"
)

func TestVerifyIp(t *testing.T) {
	// Arrange
	ctx := context.Background()
	redisConn := redis.NewRedis("localhost:6379", "Redis", 0)
	rdb, err := redisConn.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	redisRepository := repository.NewRedisRepositoryDb(rdb)

	handler := &Handler{
		IRedisRepository: redisRepository,
		Seconds:          5,
		NumberRequests:   5,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.0.1:12345"

	res := httptest.NewRecorder()

	// Act
	for i := 0; i <= handler.NumberRequests; i++ {
		verifyIp(handler, res, req)
	}

	// Assert
	require.Equal(t, http.StatusTooManyRequests, res.Code)
	expectedBody := "you have reached the maximum number of requests or actions allowed within a certain time frame\n"
	require.Equal(t, expectedBody, res.Body.String())
}

func TestVerifyToken(t *testing.T) {
	// Arrange
	ctx := context.Background()
	redisConn := redis.NewRedis("localhost:6379", "Redis", 0)
	rdb, err := redisConn.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	redisRepository := repository.NewRedisRepositoryDb(rdb)

	handler := &Handler{
		IRedisRepository: redisRepository,
		Seconds:          5,
		NumberRequests:   5,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("API_KEY", "5753bac1-fad8-490a-818e-f97074655028")

	res := httptest.NewRecorder()

	// Act
	for i := 0; i <= handler.NumberRequests; i++ {
		verifyToken(handler, res, req)
	}

	// Assert
	require.Equal(t, http.StatusTooManyRequests, res.Code)
	expectedBody := "you have reached the maximum number of requests or actions allowed within a certain time frame\n"
	require.Equal(t, expectedBody, res.Body.String())
}
