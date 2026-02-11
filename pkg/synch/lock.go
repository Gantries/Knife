// Package synch provides synchronization primitives.
//
// This file defines a distributed Lock interface for acquiring and releasing
// locks with context support and timeouts.
package synch

import (
	"context"
	"time"
)

type Lock interface {
	Lock(ctx context.Context, source, owner string, timeout time.Duration) (bool, error)
	Unlock(ctx context.Context, source, owner string) (bool, error)
}
