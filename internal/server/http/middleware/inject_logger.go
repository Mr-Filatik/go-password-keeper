// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"net/http"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// InjectLogger represents middleware for inject logger into context.
//
// Parameters:
//   - logger logging.Logger: logger.
func InjectLogger(logger logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.ToContext(r.Context(), logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
