// Package ha provides high-availability (HA) execution patterns.
//
// It supports cache-based HA executors that can register and execute
// operations with periodic callbacks for health checking and recovery.
package ha

import (
	"context"
	"time"

	"github.com/gantries/knife/pkg/cache"
	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/national"
)

type Type string

const (
	Cache Type = "cache"
)

type Properties struct {
	Type  Type             `yaml:"type" default:"redis"`
	Cache cache.Properties `yaml:"cache"`
}

type Kind string

type Executor[T any] interface {
	Register(ctxt context.Context, kind Kind, interval time.Duration, cb func(v ...T) *national.Message) (*time.Ticker,
		*national.Message)
	Exec(ctxt context.Context, kind Kind, va ...T) *national.Message
}

func New[T any](ctxt context.Context, props *Properties) (Executor[T], *national.Message) {
	switch props.Type {
	case Cache:
		return newCacheExecutor[T](ctxt, props)
	}
	return nil, errors.UnsupportedValueError.Msg("type", "ha-type", "value", props.Type)
}
