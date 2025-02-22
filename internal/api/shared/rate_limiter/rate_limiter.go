package rate_limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client      *redis.Client
	maxAttempts int
	duration    time.Duration
}

func NewRateLimiter(client *redis.Client, maxAttempts int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:      client,
		maxAttempts: maxAttempts,
		duration:    duration,
	}
}

func (rl *RateLimiter) IsAllowed(ctx context.Context, key string) (bool, error) {
	key = fmt.Sprintf("rate_limit:%s", key)

	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.client.Expire(ctx, key, rl.duration)
	}

	return count <= int64(rl.maxAttempts), nil
}
