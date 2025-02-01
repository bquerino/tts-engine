package repository

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisUsageRepository struct {
	client *redis.Client
}

func NewRedisUsageRepository(client *redis.Client) *RedisUsageRepository {
	return &RedisUsageRepository{client: client}
}

func (r *RedisUsageRepository) IncrementMessageCount(provider string) error {
	ctx := context.Background()
	_, err := r.client.Incr(ctx, provider+":messages").Result()
	return err
}

func (r *RedisUsageRepository) IncrementCharacterCount(provider string, count int) error {
	ctx := context.Background()
	_, err := r.client.IncrBy(ctx, provider+":characters", int64(count)).Result()
	return err
}

func (r *RedisUsageRepository) GetUsage(provider string) (int, int, error) {
	ctx := context.Background()

	messagesStr, err := r.client.Get(ctx, provider+":messages").Result()
	if err == redis.Nil {
		messagesStr = "0"
	} else if err != nil {
		return 0, 0, err
	}

	charactersStr, err := r.client.Get(ctx, provider+":characters").Result()
	if err == redis.Nil {
		charactersStr = "0"
	} else if err != nil {
		return 0, 0, err
	}

	messages, _ := strconv.Atoi(messagesStr)
	characters, _ := strconv.Atoi(charactersStr)

	return messages, characters, nil
}

func (r *RedisUsageRepository) ResetUsage(provider string) error {
	ctx := context.Background()
	_, err := r.client.Del(ctx, provider+":messages", provider+":characters").Result()
	return err
}
