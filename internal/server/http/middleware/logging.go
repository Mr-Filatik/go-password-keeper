// Package middleware provides functionality for HTTP middleware.
package middleware

import (
	"net/http"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/observer"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// LoggingOpts - options for logging middleware.
type LoggingOpts struct {
	EnableRequestBodyLogging  bool      // Whether to enable request body logging.
	EnableResponseBodyLogging bool      // Whether to enable response body logging.
	RouteFn                   RouteFunc // Function for forming a route.
}

// Logging represents middleware for logging HTTP handlers.
//
// Parameters:
//   - logger logging.Logger: logger;
//   - options MetricsOpts: options.
//
//nolint:funlen // the formation of log fields needs to be reworked
func Logging(logger logging.Logger, options LoggingOpts) Middleware {
	if options.RouteFn == nil {
		options.RouteFn = defaultRouteFunc()
	}

	logFn := func(
		status int,
		duration time.Duration,
		reqObs *observer.RequestObserver,
		respObs *observer.ResponseObserver,
	) {
		fields := []any{
			"duration_ms", duration.Milliseconds(),
			"request_uri", reqObs.GetURI(),
			"request_method", reqObs.GetMethod(),
			"request_path", reqObs.GetURLPath(),
			"request_query", reqObs.GetURLQuery(),
			"request_protocol", reqObs.GetProtocol(),
			"request_route", reqObs.GetRoute(),
			"request_size", reqObs.GetBodySize(),
			"response_status", status,
			"response_size", respObs.GetBodySize(),
			"request_id", reqObs.GetHeader(HeaderRequestID),
			"span_id", reqObs.GetHeader("Span-ID"),
			"trace_id", reqObs.GetHeader("Trace-ID"),
		}

		if options.EnableRequestBodyLogging {
			fields = append(fields,
				"request_body", reqObs.GetBodyString(),
			)
		}

		if options.EnableResponseBodyLogging {
			fields = append(fields,
				"response_body", respObs.GetBodyString(),
			)
		}

		logger.Info("HTTP Request-Response",
			fields...,
		)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqObs := observer.NewRequestObserver(r,
				options.EnableRequestBodyLogging, options.RouteFn)
			respObs := observer.NewResponseObserver(w, options.EnableResponseBodyLogging)

			defer func() {
				if rec := recover(); rec != nil {
					logFn(
						http.StatusInternalServerError,
						time.Since(start),
						reqObs,
						respObs,
					)

					panic(rec)
				}

				status := respObs.GetStatus()
				if status == 0 {
					status = http.StatusOK
				}

				logFn(
					status,
					time.Since(start),
					reqObs,
					respObs,
				)
			}()

			next.ServeHTTP(respObs, r)
		})
	}
}
