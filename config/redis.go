package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}
	return redis.NewClient(opt)
}
