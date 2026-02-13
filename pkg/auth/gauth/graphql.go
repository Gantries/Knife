package gauth

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gantries/knife/pkg/auth"
)

// FieldAuthenticator returns a GraphQL handler option for field-level authentication
// Deprecated: handler.Option is deprecated in gqlgen, but kept for backward compatibility
//
//lint:ignore SA1019 ignore deprecation warning for backward compatibility
func FieldAuthenticator(log *slog.Logger, whitelist Whitelist) handler.Option {
	return handler.ResolverMiddleware(
		func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			fc := graphql.GetFieldContext(ctx)
			if fc.IsMethod && !auth.IsAuthorized(ctx, fc.Field.Name) {
				var err error
				if ctx, _, err = auth.Authorize(ctx, fc.Field.Name, func(c context.Context, i *auth.Identity) *auth.Identity {
					if whitelist.In(fc.Field.Name) {
						log.Info("Field in whitelist", "field", fc.Field.Name)
						return i.Authenticated(fc.Field.Name)
					}
					return i
				}); err != nil {
					return nil, err
				}
			}
			return next(ctx)
		})
}
