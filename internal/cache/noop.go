package cache

import (
	"context"
	"time"
)

type NoOpCache struct {
	defaultTTL time.Duration
}

func NewNoOpCache(defaultTTL time.Duration) *NoOpCache {
	return &NoOpCache{
		defaultTTL: defaultTTL,
	}
}

func (c *NoOpCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrCacheMiss
}

func (c *NoOpCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (c *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *NoOpCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, ErrCacheMiss
}
