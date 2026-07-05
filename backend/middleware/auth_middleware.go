package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type contextKey string

const tokenContextKey contextKey = "jwt-token"

var ErrUnauthorized = errors.New("unauthorized: invalid token or missing")

func AuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token, ok := ctx.Value(tokenContextKey).(string)

			if !ok || !strings.HasPrefix(token, "Bearer ") {
				return nil, ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}
