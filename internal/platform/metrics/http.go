// Package metrics provides functionality for working with metrics.
package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HTTPMetrics defines the type of work with HTTP metrics of the application.
type HTTPMetrics struct {
	BaseMetrics

	requestsCounter          *prometheus.CounterVec
	requestDurationHistogram *prometheus.HistogramVec
}

// NewHTTPMetrics creates a new HTTPMetrics instance.
//
// Parameters:
//   - base BaseMetrics: a basic metric type that contains common data.
func NewHTTPMetrics(base BaseMetrics) *HTTPMetrics {
	subsystemName := "http"

	requestsCounter := base.CreateCounter(CounterOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "requests_total",
			Help:       "Total number of HTTP requests.",
			LabelNames: []string{"method", "path", "status_code"},
		},
	})

	requestDurationHistogram := base.CreateHistogram(HistogramOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "request_duration_seconds",
			Help:       "Request duration in seconds.",
			LabelNames: []string{"method", "path", "status_code"},
		},
		Buckets: nil,
	})

	return &HTTPMetrics{
		BaseMetrics:              base,
		requestsCounter:          requestsCounter,
		requestDurationHistogram: requestDurationHistogram,
	}
}

// HTTPRequestLabel describes the data required to record the metric.
type HTTPRequestLabel struct {
	// Method - HTTP request method.
	Method string

	// Path - request path template.
	//
	// You must use the template, without specifying specific data.
	// This is necessary to reduce the cardinality of metrics.
	// Correct: /path/user/{id}. Incorrect: /path/user/1, /path/user/10, etc.
	Path string

	// StatusCode - status response code.
	StatusCode int64
}

// IncRequestsCounter increments the counter by one, specifying the labels.
//
// Parameters:
//   - labels HTTPRequestLabel: labels.
func (p *HTTPMetrics) IncRequestsCounter(labels HTTPRequestLabel) {
	lbls := prometheus.Labels{
		"method":      labels.Method,
		"path":        labels.Path,
		"status_code": strconv.FormatInt(labels.StatusCode, 10),
	}

	p.requestsCounter.With(lbls).Inc()
}

// ObserveRequestDurationHistogram records the execution time for the RequestDuration histogram.
//
// Parameters:
//   - labels HTTPRequestLabel: labels;
//   - duration time.Duration: request duration.
func (p *HTTPMetrics) ObserveRequestDurationHistogram(
	labels HTTPRequestLabel,
	duration time.Duration,
) {
	lbls := prometheus.Labels{
		"method":      labels.Method,
		"path":        labels.Path,
		"status_code": strconv.FormatInt(labels.StatusCode, 10),
	}

	p.requestDurationHistogram.With(lbls).Observe(duration.Seconds())
}
