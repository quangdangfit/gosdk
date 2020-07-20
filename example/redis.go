package main

import (
	"github.com/quangdangfit/gocommon/cache/redis"
	"github.com/quangdangfit/gocommon/utils/logger"
)

func main() {
	config := redis.Config{
		Address:  "localhost:6379",
		Password: "",
		Database: 0,
	}
	option := redis.Option{
		Expiration: 100,
		KeyFn: func(key string) string {
			return "quangcache" + key
		},
	}

	cache := redis.New(config, option)
	var data interface{}

	cache.SetOrigin("quang", "quang", 100)
	cache.Get("quang", &data)
	cache.Remove("quang")
	cache.RemovePattern("qua*")

	logger.Info("data: ", data)
}
