package database

import "time"

type Config struct {
	Hosts             string
	AuthDatabase      string
	Database          string
	AuthUserName      string
	AuthPassword      string
	ConnectionTimeout int
	Env               string
	Replica           string
}

type IndexConfig struct {
	Key              []string
	Unique           bool
	DropDups         bool
	Background       bool
	Sparse           bool
	ExpireAfter      time.Duration
	Name             string
	Min, Max         int
	Minf, Maxf       float64
	BucketSize       float64
	Bits             int
	DefaultLanguage  string
	LanguageOverride string
	Weights          map[string]int
}
