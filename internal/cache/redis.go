package cache

import (
	"VoAr/internal/contextkeys"
	"VoAr/internal/logger"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewRedisCache(addr string, defaultTTL time.Duration) CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Fallback to NoOpCache if Redis is unavailable.
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return NewNoOpCache(defaultTTL)
	}

	return &RedisCache{
		client:     rdb,
		defaultTTL: defaultTTL,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	requestID := contextkeys.GetRequestID(ctx)
	cmd := r.client.Get(ctx, key)
	getValue, err := cmd.Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.CacheMiss(requestID, key)
			return nil, ErrCacheMiss
		}

		logger.CacheError(requestID, key, err)
		return nil, ErrRedisUnavailable
	}

	logger.CacheHit(requestID, key)
	return []byte(getValue), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = r.defaultTTL
	}

	requestID := contextkeys.GetRequestID(ctx)
	cmd := r.client.Set(ctx, key, value, ttl)
	err := cmd.Err()

	if err != nil {
		logger.CacheError(requestID, key, err)
		return ErrRedisUnavailable
	}

	logger.CacheSet(requestID, key)
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	requestID := contextkeys.GetRequestID(ctx)
	cmd := r.client.Del(ctx, key)
	err := cmd.Err()

	if err != nil {
		logger.CacheError(requestID, key, err)
		return ErrRedisUnavailable
	}

	logger.CacheDeleted(requestID, key)
	return nil
}

func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	requestID := contextkeys.GetRequestID(ctx)
	cmd := r.client.TTL(ctx, key)
	ttl, err := cmd.Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.CacheMiss(requestID, key)
			return 0, ErrCacheMiss
		}

		logger.CacheError(requestID, key, err)
		return 0, ErrRedisUnavailable
	}

	return ttl, nil
}
