package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache defines the methods for interacting with the cache
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

// Cacher implements the Cache interface using Redis
type Cacher struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(redisURL string) *Cacher {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &Cacher{
		client: client,
	}
}

// Get retrieves a value from the cache
func (c *Cacher) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value in the cache with a TTL
func (c *Cacher) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from the cache
func (c *Cacher) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *Cacher) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Close closes the Redis connection
func (c *Cacher) Close() error {
	return c.client.Close()
}
