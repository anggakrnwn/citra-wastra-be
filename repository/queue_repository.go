package repository

import (
	"citra-wastra-be/dto"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type QueueRepository interface {
	EnqueueXPJob(ctx context.Context, payload dto.XPJobPayload) error
	EnqueueRaw(ctx context.Context, data []byte) error
	DequeueXPJob(ctx context.Context) (string, error)
	AckJob(ctx context.Context, data string) error
}

type queueRepository struct {
	redis *redis.Client
}

func NewQueueRepository(redis *redis.Client) QueueRepository {
	return &queueRepository{redis}
}

func (r *queueRepository) EnqueueXPJob(ctx context.Context, payload dto.XPJobPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return r.redis.LPush(ctx, "xp_queue", data).Err()
}

func (r *queueRepository) EnqueueRaw(ctx context.Context, data []byte) error {
	return r.redis.LPush(ctx, "xp_queue", data).Err()
}

func (r *queueRepository) DequeueXPJob(ctx context.Context) (string, error) {
	return r.redis.LMove(ctx, "xp_queue", "xp_processing", "RIGHT", "LEFT").Result()
}

func (r *queueRepository) AckJob(ctx context.Context, data string) error {
	return r.redis.LRem(ctx, "xp_processing", 1, data).Err()
}
