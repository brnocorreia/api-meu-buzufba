package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func New(ctx context.Context, db *redis.Client) (*Cache, error) {
	return &Cache{
		redis: db,
	}, nil
}

// GetKeys receives a pattern and returns all keys that match the pattern
func (c *Cache) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fault.New("failed to get keys from cache", fault.WithError(err))
	}

	return keys, nil
}

// Delete deletes a key from the cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	err := c.redis.Del(ctx, keys...).Err()
	if err != nil {
		return fault.New("failed to delete values from cache", fault.WithError(err))
	}

	return nil
}

// GetStruct receives a key and a pointer to a struct
//
// Example:
//
//	var user dto.User
//	err := cache.GetStruct(ctx, "user:1", &user)
//	if err != nil {...}
func (c *Cache) GetStruct(ctx context.Context, key string, data any) error {
	val, err := c.get(ctx, key)
	if err != nil {
		return err // The error is already being handled in the Get function
	}

	err = json.Unmarshal(val, &data)
	if err != nil {
		return fault.New("failed to unmarshal data", fault.WithError(err))
	}

	return nil
}

func (c *Cache) SetString(ctx context.Context, key string, data string, ttl time.Duration) error {
	return c.set(ctx, key, []byte(data), ttl)
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	val, err := c.get(ctx, key)
	if err != nil {
		return "", err // The error is already being handled in the Get function
	}

	return string(val), nil
}

// SetStruct receives a key and a struct
//
// Example:
//
//	user :=  dto.User{...}
//	err := cache.SetStruct(ctx, "user:1", user, time.Minute*10)
//	if err != nil {...}
func (c *Cache) SetStruct(ctx context.Context, key string, data any, ttl time.Duration) error {
	s, err := json.Marshal(data)
	if err != nil {
		return fault.New("failed to marshal data", fault.WithError(err))
	}

	return c.set(ctx, key, s, ttl)
}

// Has checks if a key exists in the cache
func (c *Cache) Has(ctx context.Context, key string) (bool, error) {
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fault.New("failed to check if key exists in cache", fault.WithError(err))
	}

	return exists > 0, nil
}

// get is a helper function that gets a value from the cache
func (c *Cache) get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		// redis.Nil is returned when a key is not found
		if errors.Is(err, redis.Nil) {
			return nil, fault.New(
				"key not found in cache",
				fault.WithTag(fault.CACHE_MISS),
				fault.WithError(err),
			)
		}

		return nil, fault.New(
			"failed to get value from cache",
			fault.WithError(err),
		)
	}

	return val, nil
}

// set is a helper function that sets a value in the cache
func (c *Cache) set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	err := c.redis.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fault.New("failed to set value in cache", fault.WithError(err))
	}

	return nil
}

func (c *Cache) Close() error {
	return c.redis.Close()
}
