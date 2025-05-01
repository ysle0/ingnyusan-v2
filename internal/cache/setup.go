package cache

import (
	"context"
	"log"
)

// SetupCache initializes the Redis cache
func SetupCache(redisURL string) *Cacher {
	redisCache := NewRedisCache(redisURL)

	// Test Redis connection
	_, err := redisCache.Get(context.Background(), "test")
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("WARNING: Could not connect to Redis: %v", err)
		log.Printf("Continuing without caching")
		return nil
	}

	log.Printf("Connected to Redis cache at %s", redisURL)
	return redisCache
}
