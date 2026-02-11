package ha

import (
	"context"
	"time"

	"github.com/gantries/knife/pkg/cache"
	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/national"
	"github.com/gantries/knife/pkg/serde"
	"github.com/gantries/knife/pkg/synch"
)

var logger = log.New("knife/ha/redis")

const keyRedisPrefix = "knife/ha/redis/"

type cacheExecutor[T any] struct {
	redis     cache.Cache
	registry  maps.Map[Kind, func(v ...T) *national.Message]
	generator func(Kind) string
}

func newCacheExecutor[T any](ctxt context.Context, props *Properties) (Executor[T], *national.Message) {
	c, m := cache.New(ctxt, &props.Cache)
	if !m.Fine() {
		return nil, m
	}
	return &cacheExecutor[T]{
		redis:    c,
		registry: maps.Map[Kind, func(v ...T) *national.Message]{},
		generator: func(kind Kind) string {
			return keyRedisPrefix + string(kind)
		}}, errors.Yes()
}

func (e *cacheExecutor[T]) Register(ctxt context.Context, kind Kind, interval time.Duration,
	cb func(v ...T) *national.Message) (*time.Ticker, *national.Message) {
	if _, ok := e.registry[kind]; ok {
		return nil, errors.OverwriteIsForbiddenError.Msg("target", kind, "type", "ha-registry")
	}
	e.registry[kind] = cb

	ticker := time.NewTicker(interval)
	runner := synch.Runner(1, func() {
		logger.Info("Ha worker start", "kind", kind, "interval", interval)
	})
	go func() {
		for t := range ticker.C {
			_ = runner.Run()
			queue := e.generator(kind)
			ctxt := context.Background()
			l, err := e.redis.Count(ctxt, queue)
			if err != nil {
				logger.Error("Unable to count element", "error", err, "when", t, "kind", kind, "interval", interval)
			}
			if l > 0 {
				v, err := e.redis.Pop(context.Background(), queue)
				if err != nil {
					logger.Error("Unable to pop element", "error", err, "when", t, "kind", kind, "interval", interval)
				} else {
					a, _ := serde.DeserializeArray[T]([]byte(v))
					m := cb(a...)
					logger.Info("Ha result", "result", *m.Body().Raw(), "when", t, "kind", kind, "interval",
						interval)
				}
			}
		}
	}()

	return ticker, errors.Yes()
}

func (e *cacheExecutor[T]) Exec(ctxt context.Context, kind Kind, va ...T) *national.Message {
	if cb, ok := e.registry[kind]; ok {
		failed := lists.List[T]{}
		for _, v := range va {
			if m := cb(v); !m.Fine() {
				failed.Add(v)
			}
		}
		if len(failed) > 0 {
			buf, err := serde.Serialize(failed)
			if err != nil {
				logger.Error("Unable to serialize element", "error", err)
				return errors.No(err)
			}
			if err := e.redis.Push(ctxt, e.generator(kind), buf); err != nil {
				logger.Error("Unable to push element to ha queue", "error", err, "kind", kind)
				return errors.No(err)
			}
		}
		return errors.Yes()
	} else {
		return errors.NotFoundError.Build("type", "ha-callback", "value", kind)
	}
}
