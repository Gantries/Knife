package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/national"
	"github.com/gin-gonic/gin"
)

type Identity struct {
	Email             string
	Name              string
	UserName          string
	Raw               string
	authenticated     bool
	authenticatedKeys map[string]bool
}

func (i *Identity) Authenticated(id string) *Identity {
	i.authenticatedKeys[id] = true
	return i
}

const (
	keyEmail    = "email"
	keyName     = "name"
	keyUsername = "username"
)

type contextKeyType string

const HeaderIdentity = contextKeyType("x-userinfo")
const KeyIdentity = contextKeyType("authorization")

func NewIdentity(email, name, username, raw string) *Identity {
	return &Identity{
		Email:             email,
		Name:              name,
		Raw:               raw,
		UserName:          username,
		authenticated:     false,
		authenticatedKeys: make(map[string]bool),
	}
}

var logger = log.New("admin/utils/userinfo")

func parseIdentity(raw string, i *Identity) *Identity {
	if len(raw) == 0 {
		return i
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		logger.Error("Parse user info error", "error", err)
		return i
	}
	info := make(map[string]interface{})
	if err := json.Unmarshal(decodedBytes, &info); err != nil {
		logger.Error("Parse user info error", "error", err)
		return i
	}

	var email, name, username string

	if v, ok := info[keyEmail]; ok {
		if email, ok = v.(string); !ok {
			logger.Error("Missing email")
			return i
		}
	} else {
		logger.Error("Unable to parse email")
		return i
	}

	if v, ok := info[keyName]; ok {
		if name, ok = v.(string); !ok {
			logger.Error("Missing name")
			return i
		}
	} else {
		logger.Error("Unable to parse name")
		return i
	}

	if v, ok := info[keyUsername]; ok {
		if username, ok = v.(string); !ok {
			logger.Error("Missing username")
			return i
		}
	} else {
		logger.Error("Unable to parse username")
		return i
	}

	i.Email = email
	i.Name = name
	i.authenticated = true
	i.Raw = raw
	i.UserName = username
	return i
}

func withIdentity(ctx context.Context, user *Identity) context.Context {
	return context.WithValue(ctx, HeaderIdentity, user)
}

// IdentityFromContext 如果userInfoKey不存在，返回结果为nil
func IdentityFromContext(ctx context.Context) *Identity {
	if user, ok := ctx.Value(HeaderIdentity).(*Identity); ok {
		return user
	}
	return nil
}

func AuthorizationFromContext(ctx context.Context) string {
	authorization, ok := ctx.Value(KeyIdentity).(string)
	if !ok {
		return ""
	}
	return authorization
}

func Authorize(ctx context.Context, optionalId string, hook func(context.Context, *Identity) *Identity) (context.Context, *Identity, error) {
	info := parseIdentity(AuthorizationFromContext(ctx), IdentityFromContext(ctx))
	if info == nil {
		return ctx, nil, errors.MissingAuthenticationToken.LocalE(national.Tr(ctx), logger)
	}
	if hook != nil {
		info = hook(ctx, info)
	}
	if !info.authenticated {
		if y, ok := info.authenticatedKeys[optionalId]; !ok || !y {
			return ctx, nil, errors.Unauthorized.LocalE(national.Tr(ctx), logger)
		}
	}
	return context.WithValue(ctx, HeaderIdentity, info), info, nil
}

func PreAuthorize(c *gin.Context) *http.Request {
	authentication := c.Request.Header.Get(string(HeaderIdentity))
	ctx := context.WithValue(c.Request.Context(), KeyIdentity, authentication)
	return c.Request.WithContext(withIdentity(ctx, &Identity{authenticated: false, authenticatedKeys: make(map[string]bool)}))
}

func IsAuthorized(ctx context.Context, optionalId string) bool {
	if i := IdentityFromContext(ctx); i != nil {
		if i.authenticated {
			return true
		}
		if y, ok := i.authenticatedKeys[optionalId]; ok && y {
			return true
		}
	}
	return false
}
