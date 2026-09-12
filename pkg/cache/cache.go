package cache

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Cache defines the operations shared by all cache drivers.
type Cache interface {
	// Set stores the value under the key with the given ttl in seconds.
	Set(ctx context.Context, key, value string, ttl int) error
	// Get returns the value stored under the key, zero value when missing.
	Get(ctx context.Context, key string) (string, error)
	// Delete removes the key from the cache.
	Delete(ctx context.Context, key string) error
	// Exists reports whether the key exists in the cache.
	Exists(ctx context.Context, key string) (bool, error)
	// Close releases the resources held by the cache.
	Close() error
}

// Options holds the cache settings.
type Options struct {
	// Driver is the cache driver (memory/redis, case-insensitive).
	Driver string
	// TTL is the default expiration in seconds, defaulting to 300.
	TTL int
	// Addr is the redis server address in host:port form.
	Addr string
	// Password is the redis auth password, empty if auth is disabled.
	Password string
	// DB is the redis logical database number, defaulting to 0.
	DB int
}

// New creates the cache from the given options.
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
		return newRedisCache(opts, duration)
	default:
		return nil, fmt.Errorf("[cache] invalid driver %q: must be memory or redis", opts.Driver)
	}
}
