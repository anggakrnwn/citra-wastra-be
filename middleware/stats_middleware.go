package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func StatsMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path != "" {
			rdb.Incr(context.Background(), "stats:hits:"+path)
		}
		c.Next()
	}
}
