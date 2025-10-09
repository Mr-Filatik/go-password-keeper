// Package metrics provides functionality for working with metrics.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Provider represents the main entity for working with application metrics.
type Provider struct {
	// HTTP - reference to an object for working with HTTP application metrics.
	HTTP *ProviderHTTP
}

// CreateProvider creates an instance of a provider for application metrics.
//
// Parameters:
//   - namespace: common prefix for all metrics.
func CreateProvider(namespace string, appName string) *Provider {
	constLabels := prometheus.Labels{
		"app": appName,
	}

	provider := &Provider{
		HTTP: createProviderHTTP(namespace, constLabels),
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
