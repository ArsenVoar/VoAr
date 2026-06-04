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

func NewRedisCache(addr string, defaultTTL time.Duration) cache.CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return NewNoOpCache(defaultTTL)
	}

	return &RedisCache{
		client:     rdb,
		defaultTTL: defaultTTL,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := r.client.Get(ctx, key)
	getValue, err := cmd.Result()

	requestID := contextkeys.GetRequestID(ctx)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.CacheMiss(requestID, key)
			return nil, ErrCacheMiss
		} else {
			logger.CacheError(requestID, key, err)
			return nil, ErrRedisUnavailable
		}
	}

	logger.CacheHit(requestID, key)
	return []byte(getValue), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = r.defaultTTL
	}

	cmd := r.client.Set(ctx, key, value, ttl)
	err := cmd.Err()

	requestID := contextkeys.GetRequestID(ctx)

	if err != nil {
		logger.CacheError(requestID, key, err)
		return ErrRedisUnavailable
	}

	logger.CacheSet(requestID, key)
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	cmd := r.client.Del(ctx, key)
	err := cmd.Err()

	requestID := contextkeys.GetRequestID(ctx)

	if err != nil {
		logger.CacheError(requestID, key, err)
		return ErrRedisUnavailable
	}

	logger.CacheDeleted(requestID, key)
	return nil
}

func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := r.client.TTL(ctx, key)
	ttl, err := cmd.Result()

	requestID := contextkeys.GetRequestID(ctx)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.CacheMiss(requestID, key)
			return 0, ErrCacheMiss
		} else {
			logger.CacheError(requestID, key, err)
			return 0, ErrRedisUnavailable
		}
	}

	return ttl, nil
}
