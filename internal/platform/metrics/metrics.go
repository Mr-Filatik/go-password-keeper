// Package metrics provides functionality for working with metrics.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Provider represents the main entity for working with application metrics.
type Provider struct {
	// HTTP - a reference to an object for working with HTTP application metrics.
	HTTP *HTTPMetrics

	// Experiment - a reference to an object for working with experiment metrics
	Experiment *ExperimentMetrics
}

// CreateProvider creates an instance of a provider for application metrics.
//
// Parameters:
//   - namespace: common prefix for all metrics;
//   - appName: application name.
func CreateProvider(namespace string, appName string) *Provider {
	constLabels := prometheus.Labels{
		"app": appName,
	}

	baseMetrics := *NewBaseMetrics(namespace, constLabels)

	provider := &Provider{
		HTTP:       NewHTTPMetrics(baseMetrics),
		Experiment: NewExperimentMetrics(baseMetrics),
	}

	return provider
}

type handleRegister interface {
	Handle(path string, handler http.Handler)
}

// RegisterHandler registers a handler at the `/metrics` path.
//
// Parameters:
//   - router: router.
func RegisterHandler(router handleRegister) {
	router.Handle("/metrics", promhttp.Handler())
}
