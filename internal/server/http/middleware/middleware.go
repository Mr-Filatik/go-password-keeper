// Package middleware provides functionality for HTTP middleware.
package middleware

import "net/http"

// Middleware describes the type for all middlewares in an application.
type Middleware func(http.Handler) http.Handler
