package repository_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/jeanSagaz/rate-limiter/internal/domain"
	pkgRedis "github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
	pkgRepository "github.com/jeanSagaz/rate-limiter/pkg/infra/redis/repository"
	"github.com/stretchr/testify/require"
)

func TestSetRedis(t *testing.T) {
	// Arrange
	redisConn := pkgRedis.NewRedis("localhost:6379", "Redis", 0)
	ctx := context.Background()
	rdb, err := redisConn.Connect(ctx)
	if err != nil {
		log.Fatalf("error connecting to Redis")
	}

	redisRepository := pkgRepository.NewRedisRepositoryDb(rdb)

	// Act
	e := domain.Entity{Key: "test", Count: 1, Time: time.Now()}
	err = redisRepository.Set(ctx, e.Key, e, time.Second*10)

	// Assert
	require.NotEmpty(t, redisRepository)
	require.Nil(t, err)
}

func TestGetRedis(t *testing.T) {
	// Arrange
	redisConn := pkgRedis.NewRedis("localhost:6379", "Redis", 0)
	ctx := context.Background()
	rdb, err := redisConn.Connect(ctx)
	if err != nil {
		log.Fatalf("error connecting to Redis")
	}

	redisRepository := pkgRepository.NewRedisRepositoryDb(rdb)

	// Act
	e := domain.Entity{Key: "test", Count: 1, Time: time.Now()}
	errorSet := redisRepository.Set(ctx, e.Key, e, time.Second*10)

	v, err := redisRepository.Get(ctx, e.Key)

	// Assert
	require.NotEmpty(t, v)
	require.Nil(t, errorSet)
	require.Nil(t, err)
}
