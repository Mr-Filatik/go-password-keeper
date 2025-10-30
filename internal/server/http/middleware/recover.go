package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// Recover intercepts request panics and logs them.
//
// Prevents interception of http.ErrAbortHandler.
// Does not record the response for an Upgrade connection (websocket, etc.).
//
// Parameters:
//   - logger logging.Logger: logger.
func Recover(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					err, ok := rec.(error)
					if ok && errors.Is(err, http.ErrAbortHandler) {
						panic(rec)
					}

					logger.Error("HTTP Request-Response Recover", err,
						"request_id", r.Header.Get("X-Request-ID"),
						"callstack", string(debug.Stack()),
					)

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
