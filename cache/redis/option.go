package redis

import (
	"time"

	"github.com/quangdangfit/gocommon/cache"
)

type Option struct {
	KeyFn      cache.KeyFn
	Expiration time.Duration
}

func GetDefaultOption() *Option {
	return &Option{
		KeyFn:      cache.DefaultKeyFn,
		Expiration: cache.DefaultExpiration,
	}
}
