package redis

import (
	"context"
	"fmt"

	"github.com/brnocorreia/api-meu-buzufba/internal/config"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type database struct {
	db *redis.Client
}

func NewConnection(ctx context.Context, config *config.Config) (*database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logging.Error("failed to connect to redis", err,
			zap.String("journey", "redis"))
		return nil, fault.New("failed to connect to redis", fault.WithError(err))
	}

	return &database{db: client}, nil
}

func (r *database) DB() *redis.Client {
	return r.db
}

func (r *database) Close() error {
	return r.db.Close()
}
