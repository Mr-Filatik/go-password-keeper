// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"context"
	"net/http"

	logctx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/log"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// InjectLogger represents middleware for inject logger into context.
func InjectLogger(logger log.ILogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := InjectLoggerAction(r.Context(), logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// InjectLoggerAction
func InjectLoggerAction(ctx context.Context, logger log.ILogger) context.Context {
	return logctx.SetLogger(ctx, logger)
}
