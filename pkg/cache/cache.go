package cache

import (
	"context"
	"time"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/national"
)

type Cache interface {
	Ping(ctxt context.Context) error
	Push(ctx context.Context, key string, values ...interface{}) error
	Pop(ctx context.Context, key string) (string, error)
	Count(ctx context.Context, key string) (int64, error)
	Del(ctxt context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Get(ctxt context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPush(ctx context.Context, key string, fields ...string) (int64, error)
}

type Type string

type Credential struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Pool struct {
	Size                  int           `json:"size" yaml:"size" default:"60"`
	MaxRetries            int           `json:"max_retries" yaml:"max_retries" default:"3"`
	MinIdleConnections    int           `json:"min_idle_connections" yaml:"min_idle_connections" default:"10"`
	MaxIdleConnections    int           `json:"max_idle_connections" yaml:"max_idle_connections" default:"20"`
	MaxActiveConnections  int           `json:"max_active_connections" yaml:"max_active_connections" default:"30"`
	MaxConnectionIdleTime time.Duration `json:"max_connection_idle_time" yaml:"max_connection_idle_time" default:"1h"`
	MaxConnectionLifeTime time.Duration `json:"max_connection_life_time" yaml:"max_connection_life_time" default:"10h"`
}

type Properties struct {
	Type       Type               `yaml:"type" default:"redis"`
	Addresses  lists.List[string] `yaml:"addresses"`
	Database   int                `yaml:"database" default:"0"`
	Credential Credential         `yaml:"credential"`
	Pool       Pool               `yaml:"pool"`
}

const (
	Redis Type = "redis"
)

func New(ctxt context.Context, cfg *Properties) (Cache, *national.Message) {
	switch cfg.Type {
	case Redis:
		return NewRedis(ctxt, cfg)
	}
	return nil, errors.UnrecognizedError.Build("type", "cache", "value", cfg.Type)
}
