package cache

import (
	"context"
	"time"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/national"
	"github.com/redis/go-redis/v9"
)

var logger = log.New("knife/cache/redis")

type client interface {
	Ping(ctx context.Context) *redis.StatusCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	LLen(ctxt context.Context, key string) *redis.IntCmd
	SetNX(ctxt context.Context, key string, value interface{}, timeout time.Duration) *redis.BoolCmd
	Del(ctxt context.Context, keys ...string) *redis.IntCmd
	Get(ctxt context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisCache struct {
	redis client
}

func NewRedis(ctxt context.Context, cfg *Properties) (*RedisCache, *national.Message) {
	var rc RedisCache
	if cfg.Addresses.Length() > 1 {
		rc = RedisCache{redis: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           cfg.Addresses,
			Username:        cfg.Credential.Username,
			Password:        cfg.Credential.Password,
			MaxRetries:      cfg.Pool.MaxRetries,
			MinIdleConns:    cfg.Pool.MinIdleConnections,
			MaxIdleConns:    cfg.Pool.MaxIdleConnections,
			MaxActiveConns:  cfg.Pool.MaxActiveConnections,
			ConnMaxIdleTime: cfg.Pool.MaxConnectionIdleTime,
			ConnMaxLifetime: cfg.Pool.MaxConnectionLifeTime,
		})}
	} else {
		rc = RedisCache{redis: redis.NewClient(&redis.Options{
			Addr:                  cfg.Addresses[0],
			Username:              cfg.Credential.Username,
			Password:              cfg.Credential.Password,
			DB:                    cfg.Database,
			MaxRetries:            cfg.Pool.MaxRetries,
			ContextTimeoutEnabled: true,
			PoolFIFO:              false,
			PoolSize:              cfg.Pool.Size,
			MinIdleConns:          cfg.Pool.MinIdleConnections,
			MaxIdleConns:          cfg.Pool.MaxIdleConnections,
			MaxActiveConns:        cfg.Pool.MaxActiveConnections,
			ConnMaxIdleTime:       cfg.Pool.MaxConnectionIdleTime,
			ConnMaxLifetime:       cfg.Pool.MaxConnectionLifeTime,
			DisableIndentity:      false,
		})}
	}
	if err := rc.Ping(ctxt); err != nil {
		return nil, errors.No(err)
	}
	return &rc, errors.Yes()
}

func (r *RedisCache) Ping(ctxt context.Context) error {
	return r.redis.Ping(ctxt).Err()
}

func (r *RedisCache) Push(ctx context.Context, key string, values ...interface{}) error {
	return r.redis.RPush(ctx, key, values).Err()
}

func (r *RedisCache) InsertHead(ctx context.Context, key string, values ...interface{}) error {
	return r.redis.RPush(ctx, key, values).Err()
}

func (r *RedisCache) Pop(ctx context.Context, key string) (string, error) {
	return r.redis.LPop(ctx, key).Result()
}

func (r *RedisCache) Count(ctx context.Context, key string) (int64, error) {
	return r.redis.LLen(ctx, key).Result()
}

func (r *RedisCache) Lock(ctx context.Context, source, owner string, timeout time.Duration) (bool, error) {
	result, err := r.redis.SetNX(ctx, source, owner, timeout).Result()
	if result {
		return result, nil
	}
	if err != nil {
		logger.Error("Lock failed", "error", err, "source", source, "owner", owner)
		return false, nil
	}
	return false, nil
}

func (r *RedisCache) Unlock(ctx context.Context, source, owner string) (bool, error) {
	var script = redis.NewScript(`local key = KEYS[1] local value = ARGV[1] if redis.call('get',key) == value then return redis.call('del',key) else return 0 end`)
	keys := []string{source}
	values := []string{owner}
	nums, err := script.Run(ctx, r.redis.(redis.Scripter), keys, values).Bool()
	if err != nil {
		logger.Error("Unlock failed", "error", err, "source", source, "owner", owner)
		return false, nil
	}
	return nums, nil
}

func (r *RedisCache) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.redis.HSet(ctx, key, values).Result()
}

func (r *RedisCache) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return r.redis.HDel(ctx, key, fields...).Result()
}

func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	return r.redis.HGet(ctx, key, field).Result()
}

func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.redis.HGetAll(ctx, key).Result()
}

func (r *RedisCache) Del(ctxt context.Context, keys ...string) (int64, error) {
	return r.redis.Del(ctxt, keys...).Result()
}

func (r *RedisCache) RPush(ctx context.Context, key string, fields ...string) (int64, error) {
	return r.redis.RPush(ctx, key, fields).Result()
}

func (r *RedisCache) LPop(ctx context.Context, key string) (string, error) {
	return r.redis.LPop(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.redis.Set(ctx, key, value, expiration).Result()
}

func (r *RedisCache) Get(ctxt context.Context, key string) (string, error) {
	return r.redis.Get(ctxt, key).Result()
}

func (r *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.redis.Exists(ctx, keys...).Result()
}
