package main

import (
	"time"

	"github.com/quangdangfit/gosdk/cache/redis"
	"github.com/quangdangfit/gosdk/utils/logger"
)

func main() {
	config := redis.Config{
		Address:  "localhost:6379",
		Password: "",
		Database: 0,
	}
	keyFn := func(key string) string {
		return "quangcache" + key
	}

	c := redis.New(config, redis.WithExpiration(100*time.Second), redis.WithKeyFn(keyFn))
	var data interface{}

	c.SetOrigin("quang", "quang", 100*time.Second)
	c.Set("quang", "quang", 100*time.Second)
	c.Get("quang", &data)
	c.Remove("quang")
	c.RemovePattern("qua*")

	logger.Info("data: ", data)
}
