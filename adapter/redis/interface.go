package redis

import (
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// IRedis interface
type IRedis interface {
	IsConnected() bool
	Get(key string, value interface{}) error
	Set(key string, value interface{}) error
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
	Remove(keys ...string) error
	Keys(pattern string) ([]string, error)
	RemovePattern(pattern string) error
	AcquireLock(key string, value string, ttl time.Duration) (bool, error)
	AcquireLockWithRetry(key string, value string, ttl time.Duration, retryTimes int, sleepTime time.Duration) (bool, error)
	CMD() goredis.Cmdable
	Client() *goredis.Client
}
