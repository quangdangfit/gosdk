package redis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/go-redis/redis/v8"

	"github.com/quangdangfit/gosdk/cache"
	"github.com/quangdangfit/gosdk/utils/logger"
)

var ctx = context.Background()

type redis struct {
	cmd        goredis.Cmdable
	keyFn      cache.KeyFn
	expiration time.Duration
}

func New(config Config, opts ...Option) cache.Cache {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error(pong, err)
		return nil
	}

	opt := getConfig(opts...)

	return &redis{
		cmd:        rdb,
		keyFn:      opt.keyFn,
		expiration: opt.expiration,
	}
}

func (r *redis) IsConnected() bool {
	if r.cmd == nil {
		return false
	}

	_, err := r.cmd.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *redis) Get(key string, value interface{}) error {
	cacheKey := r.keyFn(key)
	return r.GetOrigin(cacheKey, value)
}
func (r *redis) GetOrigin(key string, value interface{}) error {
	strValue, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil
		}

		logger.Info("Failed to get: ", err)
		return err
	}

	logger.Infof("Get from redis %s: %s", key, strValue)

	data, _ := json.Marshal(strValue)
	err = json.Unmarshal(data, value)
	if err != nil {
		logger.Error("Failed to deserialize data", "error", err)
		return err
	}

	return nil
}

func (r *redis) Set(key string, value interface{}) error {
	cacheKey := r.keyFn(key)
	return r.SetOrigin(cacheKey, value, r.expiration)
}

func (r *redis) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	cacheKey := r.keyFn(key)
	return r.SetOrigin(cacheKey, value, expiration)
}

func (r *redis) SetOrigin(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = cache.DefaultExpiration
	}

	err := r.cmd.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Errorf("Failed to set: ", err)
		return err
	}
	logger.Infof("Set to redis %s: %s", key, value)

	return nil
}

func (r *redis) Remove(keys ...string) error {
	var cacheKeys []string
	for _, key := range keys {
		cacheKeys = append(cacheKeys, r.keyFn(key))
	}

	return r.RemoveOrigin(cacheKeys...)
}

func (r *redis) RemoveOrigin(keys ...string) error {
	err := r.cmd.Del(ctx, keys...).Err()
	if err != nil {
		logger.Errorf("Failed to delete keys %s: %s", keys, err)
		return err
	}
	logger.Info("Deleted keys: ", keys)

	return nil
}

func (r *redis) Keys(pattern string) ([]string, error) {
	keys, err := r.cmd.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Errorf("Failed to get pattern %s: %s", pattern, err)
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
		logger.Info("Not found any key with pattern: ", pattern)
		return nil
	}

	err = r.RemoveOrigin(keys...)
	if err != nil {
		return err
	}

	return nil
}
