package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GamificationRepository interface {
	IncrementUserXP(ctx context.Context, userID string, increment float64) (float64, error)
	GetUserXP(ctx context.Context, userID string) (float64, error)
}

type gamificationRepository struct {
	redis *redis.Client
}

func NewGamificationRepository(redis *redis.Client) GamificationRepository {
	return &gamificationRepository{redis}
}

func (r *gamificationRepository) IncrementUserXP(ctx context.Context, userID string, increment float64) (float64, error) {
	return r.redis.ZIncrBy(ctx, "leaderboard", increment, userID).Result()

}

func (r *gamificationRepository) GetUserXP(ctx context.Context, userID string) (float64, error) {
	return r.redis.ZScore(ctx, "leaderboard", userID).Result()
}
