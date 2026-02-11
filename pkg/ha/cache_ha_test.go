package ha

import (
	"context"
	"testing"
	"time"

	"github.com/gantries/knife/pkg/cache"
	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/national"
	"github.com/stretchr/testify/assert"
)

var properties = Properties{
	Type: Cache,
	Cache: cache.Properties{
		Type:       cache.Redis,
		Addresses:  *lists.Of[string]("127.0.0.1:6379"),
		Database:   0,
		Credential: cache.Credential{},
		Pool:       cache.Pool{},
	},
}

func TestNew(t *testing.T) {
	ctxt := context.Background()

	// Check if Redis is available before running HA tests
	testCache, m := cache.New(ctxt, &properties.Cache)
	if testCache == nil || !m.Fine() || testCache.Ping(ctxt) != nil {
		t.Skip("Redis not available, skipping HA cache tests")
		return
	}

	kind := Kind("ut")

	ha, m := New[any](ctxt, &properties)
	assert.True(t, m.Fine())
	assert.NotNil(t, ha)

	ch := make(chan string)
	go func() {
		for {
			switch s := <-ch; s {
			case "stop":
				break
			}
		}
		//ticker.Stop()
	}()
	ha.Register(ctxt, kind, 5*time.Second, func(v ...any) *national.Message {
		logger.Info("got job", "job", v)
		ch <- "stop"
		return errors.Yes()
	})
	ha.Exec(ctxt, kind, "hello", "world")
	ha.Exec(ctxt, kind, "ha job")
	logger.Info("ready to exit")
}
