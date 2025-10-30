package middleware

import (
	"net/http"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

type MetricsOption func(*MetricsOpts)

type MetricsOpts struct {
	GetRequestRouteFn func(ctx *http.Request) string // Использовать строго после next.ServeHTTP(sw, r)
}

func MetricsWithRequestRoute(reqFn func(ctx *http.Request) string) MetricsOption {
	return func(o *MetricsOpts) {
		o.GetRequestRouteFn = reqFn
	}
}

func Metrics(metricsProvider *metrics.Provider, opts ...MetricsOption) func(http.Handler) http.Handler {
	options := defaultMetricsOpts()
	applyMetricsOpts(options, opts)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: 0}

			defer func() {
				route := options.GetRequestRouteFn(r)
				duration := time.Since(start)

				if rec := recover(); rec != nil {
					metricsProvider.HTTP.IncRequestsCounter(metrics.HTTPRequestLabel{
						Method:     r.Method,
						Path:       route,
						StatusCode: http.StatusInternalServerError,
					})

					metricsProvider.HTTP.ObserveRequestDurationHistogram(metrics.HTTPRequestLabel{
						Method:     r.Method,
						Path:       route,
						StatusCode: http.StatusInternalServerError,
					}, duration)

					panic(rec)
				}

				status := sw.status
				if status == 0 {
					status = http.StatusOK
				}

				metricsProvider.HTTP.IncRequestsCounter(metrics.HTTPRequestLabel{
					Method:     r.Method,
					Path:       route,
					StatusCode: int64(status),
				})

				metricsProvider.HTTP.ObserveRequestDurationHistogram(metrics.HTTPRequestLabel{
					Method:     r.Method,
					Path:       route,
					StatusCode: int64(status),
				}, duration)
			}()

			next.ServeHTTP(sw, r)
		})
	}
}

func defaultMetricsOpts() *MetricsOpts {
	return &MetricsOpts{
		GetRequestRouteFn: func(_ *http.Request) string {
			return "none"
		},
	}
}

func applyMetricsOpts(currentOpts *MetricsOpts, opts []MetricsOption) {
	for _, apply := range opts {
		apply(currentOpts)
	}
}

// s.metricsProvider.HTTP.IncRequestsCounter(metrics.HTTPRequestLabel{
// 		Method:     r.Method,
// 		Path:       "/ping",
// 		StatusCode: http.StatusOK,
// 	})
// 	s.metricsProvider.HTTP.ObserveRequestDurationHistogram(metrics.HTTPRequestLabel{
// 		Method:     r.Method,
// 		Path:       "/ping",
// 		StatusCode: http.StatusOK,
// 	}, time.Since(start))

// Общая обёртка для захвата статуса/байтов.
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
