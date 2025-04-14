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

func TestUserSession(t *testing.T) {
	ctx := context.Background()
	id := "test-id"
	token := "test-token"
	expiration := time.Hour

	err := SetUserSession(ctx, id, token, expiration)
	assert.NoError(t, err)

	exists := UserSessionCheck(id, token)
	assert.True(t, exists)

	err = RemoveUserSession(ctx, id)
	assert.NoError(t, err)

	exists = UserSessionCheck(id, token)
	assert.False(t, exists)
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

func TestSetAdvertisement(t *testing.T) {
	ctx := context.Background()
	id := "test-id"
	expiration := time.Hour

	err := SetAdvertisement(ctx, id, expiration)
	assert.NoError(t, err)

	count, err := CountKeysByPattern(ctx, "status-wp-user-advertisement-"+id+"*")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCountKeysByPattern(t *testing.T) {
	ctx := context.Background()

	// Add some test keys
	_ = Set(ctx, "test-pattern-1", "value1", time.Hour)
	_ = Set(ctx, "test-pattern-2", "value2", time.Hour)
	_ = Set(ctx, "different-key", "value3", time.Hour)

	count, err := CountKeysByPattern(ctx, "test-pattern-*")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
