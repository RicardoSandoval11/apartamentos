package middleware

import (
	"context"
	"net/http"
)

func ExtractTokenFromHeader(ctx context.Context, r *http.Request) context.Context {
	token := r.Header.Get("Authorization")
	return context.WithValue(ctx, tokenContextKey, token)
}
