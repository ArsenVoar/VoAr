package cache

import "errors"

var ErrCacheMiss = errors.New("cache key missing")
var ErrRedisUnavailable = errors.New("redis unavailable")
