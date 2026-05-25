// Package metrics provides functionality for working with metrics.
package metrics

import (
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Provider represents the main entity for working with application metrics.
type Provider struct {
	// App - a reference to an object for working with application metrics.
	App *AppMetrics

	// HTTP - a reference to an object for working with HTTP application metrics.
	HTTP *HTTPMetrics

	// Experiment - a reference to an object for working with experiment metrics
	Experiment *ExperimentMetrics

	Log *LogMetrics
}

// CreateProvider creates an instance of a provider for application metrics.
//
// Parameters:
//   - namespace: common prefix for all metrics;
//   - appName: application name.
func CreateProvider(namespace, projectName, appName string) *Provider {
	constLabels := prometheus.Labels{
		"project": projectName,
		"app":     appName,
	}

	baseMetrics := *NewBaseMetrics(namespace, constLabels)

	provider := &Provider{
		App:        NewAppMetrics(baseMetrics),
		HTTP:       NewHTTPMetrics(baseMetrics),
		Experiment: NewExperimentMetrics(baseMetrics),
		Log:        NewLogMetrics(baseMetrics),
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
func RegisterHandler(router handleRegister) { // move to GetHandler in this and use in diag server
	router.Handle("/metrics", promhttp.Handler())
}

// LastScrapeCount - number of the last received metrics.
//
// TODO: Move to DiagnosticServer.
//
//nolint:godox,gochecknoglobals
var LastScrapeCount uint64

// RegisterHandlerWithScrapeCount registers a handler at the `/metrics` path,
// with a count of the number of metrics viewed.
//
// Parameters:
//   - router: router.
func RegisterHandlerWithScrapeCount(router handleRegister) {
	base := promhttp.Handler()

	hdlFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&LastScrapeCount, 1)
		base.ServeHTTP(w, r)
	})

	router.Handle("/metrics", hdlFn)
}
