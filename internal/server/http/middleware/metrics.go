// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"net/http"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/observer"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// MetricsOpts - options for metrics middleware.
type MetricsOpts struct {
	RouteFn RouteFunc // Function for forming a route.
}

// Metrics represents middleware for tracking HTTP metrics.
//
// Parameters:
//   - metricsProvider *metrics.Provider: metrics provider;
//   - options MetricsOpts: options.
func Metrics(metricsProvider *metrics.Provider, options MetricsOpts) Middleware {
	if options.RouteFn == nil {
		options.RouteFn = defaultRouteFunc()
	}

	metrFn := func(duration time.Duration, route string, method string, statusCode int) {
		metricsProvider.HTTP.IncRequestsCounter(metrics.HTTPRequestLabel{
			Method:     method,
			Path:       route,
			StatusCode: int64(statusCode),
		})

		metricsProvider.HTTP.ObserveRequestDurationHistogram(metrics.HTTPRequestLabel{
			Method:     method,
			Path:       route,
			StatusCode: int64(statusCode),
		}, duration)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqObs := observer.NewRequestObserver(r, false, options.RouteFn)
			respObs := observer.NewResponseObserver(w, false)

			defer func() {
				duration := time.Since(start)

				if rec := recover(); rec != nil {
					metrFn(
						duration,
						reqObs.GetRoute(),
						reqObs.GetMethod(),
						http.StatusInternalServerError)

					panic(rec)
				}

				status := respObs.GetStatus()
				if status == 0 {
					status = http.StatusOK
				}

				metrFn(duration, reqObs.GetRoute(), reqObs.GetMethod(), status)
			}()

			next.ServeHTTP(respObs, r)
		})
	}
}
