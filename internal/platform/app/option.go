package app

import (
	"github.com/mr-filatik/go-password-keeper/internal/platform/http/diagnostic"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// Option describes optional settings for the App.
type Option func(*App)

// WithStopLaunchingOnError stops other components from launching when the first error is received.
// Only affects the launch of a SequentialComponent.
func WithStopLaunchingOnError() Option {
	return func(a *App) {
		a.stopLaunchingOnError = true
	}
}

// WithStartStopMetrics sets the provider for recording metrics about the start and stop of application components.
func WithStartStopMetrics(metrProv IAppMetrics) Option {
	return func(a *App) {
		a.metricsProvider = metrProv
	}
}

func WithDiagnosticServer(addr string, metr *metrics.Provider) Option {
	return func(a *App) {
		a.diagServer = diagnostic.NewServer(diagnostic.ServerConfig{
			Address:         addr,
			MetricsProvider: metr,
		}, a.logger)
	}
}
