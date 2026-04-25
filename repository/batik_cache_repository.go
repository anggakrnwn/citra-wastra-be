package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type BatikCacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type batikCacheRepository struct {
	redis *redis.Client
}

func NewBatikCacheRepository(redis *redis.Client) BatikCacheRepository {
	return &batikCacheRepository{redis}
}

func (r *batikCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, key, jsonData, expiration).Err()
}

func (r *batikCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}
