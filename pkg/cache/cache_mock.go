package cache

import (
	"context"
	"time"
)

type CacheMock struct {
}

func (c CacheMock) Ping(ctxt context.Context) error {
	return nil
}
func (c CacheMock) Push(ctx context.Context, key string, values ...interface{}) error {
	return nil
}
func (c CacheMock) Pop(ctx context.Context, key string) (string, error) {
	return "mock", nil
}
func (c CacheMock) Count(ctx context.Context, key string) (int64, error) {
	return 1, nil
}
func (c CacheMock) Del(ctxt context.Context, keys ...string) (int64, error) {
	return 1, nil
}
func (c CacheMock) Exists(ctx context.Context, keys ...string) (int64, error) {
	return 1, nil
}
func (c CacheMock) Get(ctxt context.Context, key string) (string, error) {
	return "mock", nil
}
func (c CacheMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return "", nil
}
func (c CacheMock) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return 1, nil
}
func (c CacheMock) HGet(ctx context.Context, key, field string) (string, error) {
	return "mock", nil
}
func (c CacheMock) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return map[string]string{"mock": "true"}, nil
}
func (c CacheMock) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 1, nil
}
func (c CacheMock) LPop(ctx context.Context, key string) (string, error) {
	return "mock", nil
}
func (c CacheMock) RPush(ctx context.Context, key string, fields ...string) (int64, error) {
	return 1, nil
}
