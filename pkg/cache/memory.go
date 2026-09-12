package cache

import (
	"context"
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// memoryCache is the in-memory cache implementation.
type memoryCache struct {
	client *gocache.Cache
	ttl    time.Duration
}

// newMemoryCache creates an in-memory cache.
func newMemoryCache(ttl time.Duration) *memoryCache {
	return &memoryCache{client: gocache.New(ttl, 2*ttl), ttl: ttl}
}

// Set stores the value under the key with the given ttl in seconds.
func (c *memoryCache) Set(ctx context.Context, key, value string, ttl int) error {
	// Use a single default so both drivers share the same expiration semantics.
	if ttl <= 0 {
		ttl = int(c.ttl.Seconds())
	}
	c.client.Set(key, value, time.Duration(ttl)*time.Second)
	return nil
}

// Get returns the value stored under the key, zero value when missing.
func (c *memoryCache) Get(ctx context.Context, key string) (string, error) {
	value, found := c.client.Get(key)
	if !found {
		return "", nil
	}
	s, ok := value.(string)
	if !ok {
		// Set only writes strings, so this guards against future misuse.
		return "", fmt.Errorf("[cache] key %s holds an unexpected value type", key)
	}
	return s, nil
}

// Delete removes the key from the cache.
func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.client.Delete(key)
	return nil
}

// Exists reports whether the key exists in the cache.
func (c *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := c.client.Get(key)
	return found, nil
}

// Close releases the resources held by the cache.
func (c *memoryCache) Close() error {
	return nil
}
