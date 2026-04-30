package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GamificationRepository interface {
	IncrementUserXP(ctx context.Context, userID string, increment float64) (float64, error)
	GetUserXP(ctx context.Context, userID string) (float64, error)
	SetUserXP(ctx context.Context, userID string, xp float64) error
	GetLeaderboard(ctx context.Context, limit int) ([]redis.Z, error)
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

func (r *gamificationRepository) SetUserXP(ctx context.Context, userID string, xp float64) error {
	key := "leaderboard"

	err := r.redis.ZAdd(ctx, key, redis.Z{
		Score:  xp,
		Member: userID,
	}).Err()

	return err
}

func (r *gamificationRepository) GetLeaderboard(ctx context.Context, limit int) ([]redis.Z, error) {
	return r.redis.ZRevRangeWithScores(ctx, "leaderboard", 0, int64(limit-1)).Result()
}
