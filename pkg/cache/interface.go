package cache

import "time"

type RedisCacheService interface {
	Get(key string, dest any) error
	Set(key string, value any, ttl time.Duration) error
	Exits(key string) (bool, error)
	Clear(key string) error
}