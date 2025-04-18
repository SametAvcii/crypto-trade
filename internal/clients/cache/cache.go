package cache

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
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
			log.Println("Error connecting to Redis: ", err.Error())
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

func RedisAlive(ctx context.Context, rds config.Redis) {

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

func Set(ctx context.Context, key string, value interface{}, ex_time time.Duration) error {
	return client.Set(ctx, key, value, ex_time).Err()
}
