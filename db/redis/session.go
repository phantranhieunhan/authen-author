package redis

import (
	"context"
	"errors"
	"time"

	"github.com/phantranhieunhan/authen-author/adapter/redis"
	goredis "github.com/redis/go-redis/v9"
)

var (
	ErrNilSession = errors.New("error nil session")
)

type SessionStore interface {
	CreateToken(ctx context.Context, key string, expiredTime time.Duration) error
	DeleteToken(ctx context.Context, key string) error
	GetToken(ctx context.Context, key string) (string, error)
}

type RedisStore struct {
	redis redis.IRedis
}

func NewStore(db redis.IRedis) *RedisStore {
	return &RedisStore{
		redis: db,
	}
}

func (r *RedisStore) CreateToken(ctx context.Context, key string, expiredTime time.Duration) error {
	err := r.redis.SetWithExpiration(key, true, expiredTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) DeleteToken(ctx context.Context, key string) error {
	err := r.redis.Remove(key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) GetToken(ctx context.Context, key string) (string, error) {
	var value string
	err := r.redis.Get(key, &value)
	if err != nil {
		if err == goredis.Nil {
			return "", ErrNilSession
		}
		return "", err
	}

	return value, nil
}
