// Package middleware provides functionality for HTTP middleware.
package middleware

import "net/http"

// Middleware describes the type for all middlewares in an application.
type Middleware func(http.Handler) http.Handler

// RouteFunc describes the type of function for getting a route.
type RouteFunc func(*http.Request) string

func defaultRouteFunc() RouteFunc {
	return func(_ *http.Request) string {
		return "unknown"
	}
}
