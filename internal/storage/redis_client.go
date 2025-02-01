package storage

import (
	"context"
	"os"
	"time"

	"tts-engine/internal/monitoring"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	monitoring.InfoLog("🔄 Initializing Redis connection...", nil)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:            os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password:        "",
		DB:              0,
		PoolSize:        10,
		MinIdleConns:    3,
		ConnMaxLifetime: 5 * time.Minute,
	})

	// Connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		monitoring.ErrorLog("Failed to connect to Redis", err, nil)
		return
	}

	monitoring.InfoLog("Successfully connected to Redis!", nil)
}
