package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/nitishm/go-rejson/v4"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	goredis "github.com/redis/go-redis/v9"
)

type redis struct {
	cmd    goredis.Cmdable
	client *goredis.Client
	rejson *rejson.Handler
	cfg    *config
}

// New Redis interface with config
func New(opts ...Option) (IRedis, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(context.Background(), rdb)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)

	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error(pong, err)
		return nil, err
	}

	return &redis{
		cmd:    rdb,
		cfg:    cfg,
		client: rdb,
		rejson: rh,
	}, nil
}

func (r *redis) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	if r.cmd == nil {
		return false
	}

	_, err := r.cmd.Ping(ctx).Result()
	return err == nil
}

func (r *redis) Get(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	vType := reflect.TypeOf(value)
	if vType.Kind() != reflect.Pointer {
		return fmt.Errorf("redis.Get value should be a Pointer but %T", value)
	}

	var err error
	elem := reflect.ValueOf(value).Elem()
	kind := elem.Kind()
	switch kind {
	case reflect.Bool:
		var boolVal bool
		boolVal, err = r.cmd.Get(ctx, key).Bool()
		elem.SetBool(boolVal)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		var val int
		val, err = r.cmd.Get(ctx, key).Int()
		elem.SetInt(int64(val))
	case reflect.Float32, reflect.Float64:
		var val float64
		val, err = r.cmd.Get(ctx, key).Float64()
		elem.SetFloat(val)
	default:
		strValue, err := r.cmd.Get(ctx, key).Result()

		if err != nil {
			return err
		}

		b, err := json.Marshal(&strValue)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &value)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *redis) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	err := r.cmd.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	err := r.cmd.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Remove(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	err := r.cmd.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Keys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	keys, err := r.cmd.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *redis) RemovePattern(pattern string) error {
	keys, err := r.Keys(pattern)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.Remove(keys...)
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) AcquireLock(key string, value string, ttl time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Timeout)*time.Second)
	defer cancel()

	// Set the lock value using SET command with NX and PX options
	result, err := r.cmd.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *redis) AcquireLockWithRetry(key string, value string, ttl time.Duration, retryTimes int, sleepTime time.Duration) (bool, error) {
	retryCount := 0

	for {
		if retryCount >= retryTimes {
			return false, nil
		}
		retryCount += 1
		// Try to set the lock value using SET command with NX and PX options
		result, err := r.cmd.SetNX(context.Background(), key, value, ttl).Result()
		if err != nil {
			return false, err
		}

		if result {
			// Lock successfully acquired
			return true, nil
		}
		// Lock already held by another client, wait and retry
		time.Sleep(sleepTime)
	}
}

func (r *redis) CMD() goredis.Cmdable {
	return r.cmd
}

func (r *redis) Client() *goredis.Client {
	return r.client
}
