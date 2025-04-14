package cache

import (
	"context"
	"testing"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInitRedis(t *testing.T) {
	rds := config.Redis{
		Host: "localhost",
		Port: "6379",
		Pass: "",
	}

	InitRedis(rds)
	assert.NotNil(t, client)
}

func TestRedisClient(t *testing.T) {
	rds := config.Redis{
		Host: "localhost",
		Port: "6379",
		Pass: "",
	}
	InitRedis(rds)

	c := RedisClient()
	assert.NotNil(t, c)
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := time.Hour

	err := Set(ctx, key, value, expiration)
	assert.NoError(t, err)

	val, err := client.Get(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, value, val)
}
