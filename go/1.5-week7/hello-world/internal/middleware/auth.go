package middleware

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	kratosmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
)

// Server returns a server middleware that extracts user id from Authorization.
// Supported token forms:
//   - Bearer login-token-123
//   - login-token-123
//   - 123
func Server() kratosmiddleware.Middleware {
	return func(next kratosmiddleware.Handler) kratosmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}

			token := extractToken(tr.RequestHeader().Get("Authorization"))
			if token == "" {
				return next(ctx, req)
			}

			userID, err := parseUserID(token)
			if err != nil {
				return nil, err
			}

			md := metadata.Metadata{}
			if old, ok := metadata.FromServerContext(ctx); ok {
				md = metadata.New(old)
			}
			idStr := strconv.FormatInt(userID, 10)
			md["user_id"] = []string{idStr}
			md["x-user-id"] = []string{idStr}
			ctx = metadata.NewServerContext(ctx, md)
			return next(ctx, req)
		}
	}
}

func extractToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) >= 7 && strings.EqualFold(value[:7], "Bearer ") {
		return strings.TrimSpace(value[7:])
	}
	return value
}

func parseUserID(token string) (int64, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, errors.Unauthorized("UNAUTHORIZED", "missing token")
	}

	// try plain numeric id
	if id, err := strconv.ParseInt(token, 10, 64); err == nil && id > 0 {
		return id, nil
	}

	// try suffix pattern (login-token-123)
	if idx := strings.LastIndex(token, "-"); idx >= 0 && idx < len(token)-1 {
		if id, err := strconv.ParseInt(token[idx+1:], 10, 64); err == nil && id > 0 {
			return id, nil
		}
	}

	// try JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// only allow HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("UNAUTHORIZED", "unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return 0, errors.Unauthorized("UNAUTHORIZED", "invalid token")
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		// try common claim keys
		for _, k := range []string{"user_id", "uid", "sub", "id"} {
			if v, ok := claims[k]; ok {
				switch val := v.(type) {
				case float64:
					if int64(val) > 0 {
						return int64(val), nil
					}
				case string:
					if id, err := strconv.ParseInt(val, 10, 64); err == nil && id > 0 {
						return id, nil
					}
				}
			}
		}
	}

	return 0, errors.Unauthorized("UNAUTHORIZED", "invalid token")
}
