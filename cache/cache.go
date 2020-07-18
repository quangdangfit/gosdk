package cache

import (
	"time"
)

const (
	DefaultExpiration = 24 * time.Hour
)

type Cache interface {
	IsConnected() bool
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Remove(keys ...string) error
	GetOrigin(key string, value interface{}) error
	SetOrigin(key string, value interface{}, expiration time.Duration) error
	RemoveOrigin(keys ...string) error
	RemovePattern(pattern string) error
	Keys(pattern string) ([]string, error)
}

// KeyFn defines a transformer for cache keys
type KeyFn func(string) string

func DefaultKeyFn(key string) string {
	return key
}
