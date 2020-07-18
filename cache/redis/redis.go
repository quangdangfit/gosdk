package redis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/go-redis/redis/v8"

	"gitlab.com/quangdangfit/gocommon/cache"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

var ctx = context.Background()

type redis struct {
	cmd    goredis.Cmdable
	option Option
}

func New(config Config, option Option) cache.Cache {
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

	return &redis{
		cmd:    rdb,
		option: option,
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
	cacheKey := r.option.KeyFn(key)
	return r.GetOrigin(cacheKey, value)
}

func (r *redis) Set(key string, value interface{}, expiration time.Duration) error {
	cacheKey := r.option.KeyFn(key)
	return r.SetOrigin(cacheKey, value, expiration)
}

func (r *redis) Remove(keys ...string) error {
	var cacheKeys []string
	for _, key := range keys {
		cacheKeys = append(cacheKeys, r.option.KeyFn(key))
	}

	return r.RemoveOrigin(cacheKeys...)
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
	err = json.Unmarshal(data, &value)
	if err != nil {
		logger.Error("Failed to deserialize data", "error", err)
		return err
	}

	return nil
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
