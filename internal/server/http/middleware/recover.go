// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// Recover intercepts request panics and logs them.
//
// Prevents interception of http.ErrAbortHandler.
// Does not record the response for an Upgrade connection (websocket, etc.).
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					err, ok := rec.(error)
					if ok && errors.Is(err, http.ErrAbortHandler) {
						panic(rec)
					}

					log.CtxError(r.Context(), "HTTP Request-Response Recover", err,
						log.WithStackTraceField(string(debug.Stack())))
					//"request_id", r.Header.Get(HeaderRequestID),

					if strings.EqualFold(r.Header.Get("Connection"), "Upgrade") {
						return
					}

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
