package middleware

import (
	"net/http"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/observer"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

type LoggingOption func(*LoggingOpts)

type LoggingOpts struct {
	EnableBodyRequestLogging  bool
	EnableBodyResponseLogging bool
	//MaskBody                  func(data []byte) []byte
	//MaskURL                   func(uri *url.URL)
	GetRequestRouteFn func(ctx *http.Request) string // Использовать строго после next.ServeHTTP(sw, r)
}

func LoggingWithEnableBodyLogging() LoggingOption {
	return func(o *LoggingOpts) {
		o.EnableBodyRequestLogging = true
		o.EnableBodyResponseLogging = true
	}
}

func LoggingWithRequestRoute(reqFn func(ctx *http.Request) string) LoggingOption {
	return func(o *LoggingOpts) {
		o.GetRequestRouteFn = reqFn
	}
}

// func LoggingWithMaskBody(fn func(data []byte) []byte) LoggingOption {
// 	return func(o *LoggingOpts) {
// 		o.MaskBody = fn
// 	}
// }

// func LoggingWithMaskURL(fn func(uri *url.URL)) LoggingOption {
// 	return func(o *LoggingOpts) {
// 		o.MaskURL = fn
// 	}
// }

func Logging(logger logging.Logger, opts ...LoggingOption) func(http.Handler) http.Handler {
	options := defaultLoggingOpts()
	applyLoggingOpts(options, opts)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqObs := observer.NewRequestObserver(r, options.EnableBodyRequestLogging, options.GetRequestRouteFn)
			respObs := observer.NewResponseObserver(w, options.EnableBodyResponseLogging)

			start := time.Now()

			defer func() {
				if rec := recover(); rec != nil {
					fields := []any{
						"duration_ms", time.Since(start).Milliseconds(),
						"request_uri", reqObs.GetURI(),
						"request_method", reqObs.GetMethod(),
						"request_path", reqObs.GetURLPath(),
						"request_query", reqObs.GetURLQuery(),
						"request_protocol", reqObs.GetProtocol(),
						"request_route", reqObs.GetRoute(),
						"request_size", reqObs.GetBodySize(),
						"response_status", http.StatusInternalServerError,
						"response_size", respObs.GetBodySize(),
						"request_id", r.Header.Get("X-Request-ID"),
						"span_id", "",
						"trace_id", "",
					}

					if options.EnableBodyRequestLogging {
						fields = append(fields,
							"request_body", reqObs.GetBodyString(),
						)
					}

					if options.EnableBodyResponseLogging {
						fields = append(fields,
							"response_body", respObs.GetBodyString(),
						)
					}

					logger.Info("HTTP Request-Response",
						fields...,
					)

					panic(rec)
				}

				status := respObs.GetStatus()
				if status == 0 {
					status = http.StatusOK
				}

				fields := []any{
					"duration_ms", time.Since(start).Milliseconds(),
					"request_uri", reqObs.GetURI(),
					"request_method", reqObs.GetMethod(),
					"request_path", reqObs.GetURLPath(),
					"request_query", reqObs.GetURLQuery(),
					"request_protocol", reqObs.GetProtocol(),
					"request_route", reqObs.GetRoute(),
					"request_size", reqObs.GetBodySize(),
					"response_status", status,
					"response_size", respObs.GetBodySize(),
					"request_id", r.Header.Get("X-Request-ID"),
					"span_id", "",
					"trace_id", "",
				}

				if options.EnableBodyRequestLogging {
					fields = append(fields,
						"request_body", reqObs.GetBodyString(),
					)
				}

				if options.EnableBodyResponseLogging {
					fields = append(fields,
						"response_body", respObs.GetBodyString(),
					)
				}

				logger.Info("HTTP Request-Response",
					fields...,
				)
			}()

			next.ServeHTTP(respObs, r)
		})
	}
}

func defaultLoggingOpts() *LoggingOpts {
	return &LoggingOpts{
		EnableBodyRequestLogging:  false,
		EnableBodyResponseLogging: false,
		GetRequestRouteFn: func(r *http.Request) string {
			return ""
		},
		// MaskBody: func(data []byte) []byte {
		// 	return data
		// },
		// MaskURL: func(_ *url.URL) {},
	}
}

func applyLoggingOpts(currentOpts *LoggingOpts, opts []LoggingOption) {
	for _, apply := range opts {
		apply(currentOpts)
	}
}
