package app

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
		a.metrProv = metrProv
	}
}
