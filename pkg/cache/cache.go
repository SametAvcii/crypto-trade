package cache

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	client      *redis.Client
	client_once sync.Once
)

func InitRedis(rds config.Redis) {
	log.Println("Redis connection started...")
	client_once.Do(func() {
		rdc := redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(rds.Host, rds.Port),
			Password: rds.Pass,
			DB:       0, // use default DB
		})
		// -----> control
		var ctx = context.Background()
		_, err := rdc.Ping(ctx).Result()
		if err != nil {
			log.Panic("Error connecting to Redis: ", err.Error())
		}
		client = rdc
	})
}

func RedisClient() *redis.Client {
	if client == nil {
		log.Panic("Redis client is not initialized. Call InitRedis first.")
	}
	return client
}

func RedisAlive(rds config.Redis) {

	timeTicker := time.NewTicker(15 * time.Second)

	for range timeTicker.C {
		rdc := redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(rds.Host, rds.Port),
			Password: rds.Pass,
			DB:       0, // use default DB
		})
		var ctx = context.Background()
		_, err := rdc.Ping(ctx).Result()
		if err != nil {
			log.Println("Error connecting to Redis: ", err.Error())
			continue
		}
		client = rdc
	}

}

func UserSessionCheck(id string, token string) bool {
	val, err := client.Get(context.Background(), fmt.Sprintf("status-wp-user-auth-token-%v", id)).Result()
	if err != nil {
		return false
	}
	return val == token
}

func SetUserSession(ctx context.Context, id string, token string, ex_time time.Duration) error {
	return client.Set(ctx, fmt.Sprintf("status-wp-user-auth-token-%v", id), token, ex_time).Err()
}

func Set(ctx context.Context, key string, value interface{}, ex_time time.Duration) error {
	return client.Set(ctx, key, value, ex_time).Err()
}

func RemoveUserSession(ctx context.Context, id string) error {
	return client.Del(ctx, fmt.Sprintf("status-wp-user-auth-token-%v", id)).Err()
}

func SetAdvertisement(ctx context.Context, id string, ex_time time.Duration) error {
	random_uuid := uuid.New().String()
	return client.Set(ctx, fmt.Sprintf("status-wp-user-advertisement-%v-%v", id, random_uuid), "", ex_time).Err()
}

func CountKeysByPattern(ctx context.Context, pattern string) (int, error) {
	var (
		cursor     uint64
		total_keys int
		for_secure int = 100
	)

	for {
		keys, newCursor, err := client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return 0, err
		}

		total_keys += len(keys)
		cursor = newCursor

		if cursor == 0 || for_secure == 0 {
			break
		}
		for_secure--
	}

	return total_keys, nil
}
