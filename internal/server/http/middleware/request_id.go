// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mr-filatik/go-password-keeper/internal/platform/context"
)

// HeaderRequestID is the name of the HTTP header "X-Request-Id".
const HeaderRequestID = "X-Request-Id"

// CtxKeyXRequestID - key for the HTTP header "X-Request-ID".
//
//nolint:gochecknoglobals // Migrate from platform to http
var CtxKeyXRequestID = &context.CtxKey{Name: "x-request-id"}

// RequestID creates a middleware for setting the request ID.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestID := r.Header.Get(HeaderRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
				r.Header.Set(HeaderRequestID, requestID)
			}

			ctx = context.WithValue(ctx, CtxKeyXRequestID, requestID)

			w.Header().Set(HeaderRequestID, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
