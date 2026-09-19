package cache

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Cache defines the operations of a cache driver.
type Cache interface {
	// Set stores the value under the key.
	Set(ctx context.Context, key, value string, ttl int) error
	// Get loads the value under the key, empty when missing.
	Get(ctx context.Context, key string) (string, error)
	// Delete removes the key.
	Delete(ctx context.Context, key string) error
	// Exists reports whether the key exists.
	Exists(ctx context.Context, key string) (bool, error)
	// Close releases the cache resources.
	Close() error
}

// Options holds the cache settings.
type Options struct {
	Driver   string // cache driver (memory/redis)
	TTL      int    // expiration in seconds (default to 300)
	Addr     string // redis address (host:port)
	Password string // redis password
	DB       int    // redis database (default to 0)
}

func New(opts Options) (Cache, error) {
	if opts.Driver == "" {
		return nil, fmt.Errorf("[cache] driver is required")
	}

	ttl := opts.TTL
	if ttl <= 0 {
		ttl = 300
	}
	duration := time.Duration(ttl) * time.Second

	switch strings.ToLower(opts.Driver) {
	case "memory":
		return newMemoryCache(duration), nil
	case "redis":
		if opts.Addr == "" {
			return nil, fmt.Errorf("[cache] addr is required for redis")
		}
		c, err := newRedisCache(opts, duration)
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, fmt.Errorf("[cache] invalid driver %q: must be memory or redis", opts.Driver)
	}
}
