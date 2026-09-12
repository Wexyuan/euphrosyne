package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache is the redis-backed cache implementation.
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// newRedisCache creates a redis-backed cache.
func newRedisCache(opts Options, ttl time.Duration) (*redisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
	// Ping once so a bad address fails at startup instead of on first use.
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("[cache] ping redis error: %w", err)
	}
	return &redisCache{client: client, ttl: ttl}, nil
}

// Set stores the value under the key with the given ttl in seconds.
func (c *redisCache) Set(ctx context.Context, key, value string, ttl int) error {
	// Use a single default so both drivers share the same expiration semantics.
	if ttl <= 0 {
		ttl = int(c.ttl.Seconds())
	}
	return c.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

// Get returns the value stored under the key, zero value when missing.
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("[cache] get key %s error: %w", key, err)
	}
	return value, nil
}

// Delete removes the key from the cache.
func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists reports whether the key exists in the cache.
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("[cache] check key %s exists error: %w", key, err)
	}
	return n > 0, nil
}

// Close releases the resources held by the cache.
func (c *redisCache) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("[cache] close redis client error: %w", err)
	}
	return nil
}
