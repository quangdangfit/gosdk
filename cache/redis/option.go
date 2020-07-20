package redis

import (
	"time"

	"github.com/quangdangfit/gosdk/cache"
)

type Option interface {
	apply(*option)
}

type option struct {
	keyFn      cache.KeyFn
	expiration time.Duration
}

type optionFn func(*option)

func (optFn optionFn) apply(opt *option) {
	optFn(opt)
}

func WithKeyFn(fn cache.KeyFn) Option {
	return optionFn(func(opt *option) {
		opt.keyFn = fn
	})
}

func WithExpiration(exp time.Duration) Option {
	return optionFn(func(opt *option) {
		opt.expiration = exp
	})
}

func getConfig(opts ...Option) *option {
	conf := option{
		keyFn:      cache.DefaultKeyFn,
		expiration: cache.DefaultExpiration,
	}

	for _, opt := range opts {
		opt.apply(&conf)
	}

	return &conf
}
