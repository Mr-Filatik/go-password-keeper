// Package metrics provides functionality for working with metrics.
package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ProviderHTTP is an object for working with HTTP application metrics.
type ProviderHTTP struct {
	requestsTotal   prometheus.CounterVec
	requestDuration prometheus.HistogramVec
}

const subsystemHTTP = "http"

func createProviderHTTP(namespace string, constLabels prometheus.Labels) *ProviderHTTP {
	provider := &ProviderHTTP{
		requestsTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				ConstLabels: constLabels,
				Namespace:   namespace,
				Subsystem:   subsystemHTTP,
				Name:        "requests_total",
				Help:        "Total number of HTTP requests.",
			},
			HTTPRequestLabelNames(),
		),
		requestDuration: *prometheus.NewHistogramVec(
			//nolint:exhaustruct // Native Histogram is not used
			prometheus.HistogramOpts{
				ConstLabels: constLabels,
				Namespace:   namespace,
				Subsystem:   subsystemHTTP,
				Name:        "request_duration_seconds",
				Help:        "Request duration in seconds.",
				Buckets:     prometheus.DefBuckets, // or []float64{0.05,0.1,0.2,0.3,0.5,0.75,1,2}
			},
			HTTPRequestLabelNames(),
		),
	}

	prometheus.MustRegister(
		provider.requestsTotal,
		provider.requestDuration,
	)

	return provider
}

// HTTPRequestLabel describes common labels for all requests.
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

// HTTPRequestLabelNames returns common label names for all requests.
func HTTPRequestLabelNames() []string {
	return []string{"method", "path", "status_code"}
}

// IncRequestsTotal increments the RequestsTotal counter by one.
func (p *ProviderHTTP) IncRequestsTotal(labels HTTPRequestLabel) {
	lbls := prometheus.Labels{
		"method":      labels.Method,
		"path":        labels.Path,
		"status_code": strconv.FormatInt(labels.StatusCode, 10),
	}

	p.requestsTotal.With(lbls).Inc()
}

// ObserveRequestDuration records the execution time for the RequestDuration histogram.
func (p *ProviderHTTP) ObserveRequestDuration(
	labels HTTPRequestLabel,
	duration time.Duration,
) {
	lbls := prometheus.Labels{
		"method":      labels.Method,
		"path":        labels.Path,
		"status_code": strconv.FormatInt(labels.StatusCode, 10),
	}

	p.requestDuration.With(lbls).Observe(duration.Seconds())
}
