package mock

import (
	"context"

	"github.com/gantries/knife/pkg/auth"
)

func WithIdentity(ctx context.Context, user *auth.Identity) context.Context {
	return context.WithValue(ctx, auth.HeaderIdentity, user)
}
