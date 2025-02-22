package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	REDIS_URL = "REDIS_URL"
)

func NewRedisConnection(
	ctx context.Context,
) (*redis.Client, error) {
	redis_uri := os.Getenv(REDIS_URL)

	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
