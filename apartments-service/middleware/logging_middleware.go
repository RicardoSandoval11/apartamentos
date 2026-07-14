package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-kit/kit/endpoint"
)

func LoggingMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			begin := time.Now()

			response, err := next(ctx, request)

			level := slog.LevelInfo
			if err != nil {
				level = slog.LevelError
			}

			slog.Log(
				ctx,
				level,
				"Endpoint completed",
				"Duration",
				time.Since(begin),
				"err",
				err,
			)

			return response, err
		}
	}
}
