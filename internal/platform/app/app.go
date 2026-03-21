// Package app provides everything needed to work with application components.
package app

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/diagnostic"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// App provides a framework for starting and stopping application containers.
type App struct {
	logger          log.ILogger
	mainComponent   IComponent
	metricsProvider IAppMetrics
	diagServer      *diagnostic.Server

	stopLaunchingOnError bool // whether to stop the launch when the first error is received
}

// New creates an instance of the App structure.
func New(logger log.ILogger, diagServer *diagnostic.Server, opts ...Option) *App {
	app := &App{
		logger:               logger,
		mainComponent:        &nopComponent{logger},
		metricsProvider:      nil,
		diagServer:           diagServer,
		stopLaunchingOnError: false,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// RegisterComponents registers the specified components.
//
// By default, all components are registered as running sequentially.
// If you need to run components in parallel, you must use
// the RegisterMainComponent function together with ParallelComponent.
func (a *App) RegisterComponents(components ...IComponent) {
	a.mainComponent = a.WithSequentialComponent(components...)
}

// RegisterMainComponent registers the main component.
//
// The main component can be:
//   - SequentialComponent to run components sequentially;
//   - ParallelComponent to run components in parallel;
//   - any other type implementing the IComponent interface.
func (a *App) RegisterMainComponent(mainComponent IComponent) {
	a.mainComponent = mainComponent
}

// Start begins launching all application components, starting with the main component.
//
// The error is already logged within the function; the error can be used to terminate the application.
// This may need to be reconsidered in the future; for now, it's noted in the comment.
func (a *App) Start(ctx context.Context) error {
	startTime := time.Now().UTC()

	startErr := a.mainComponent.Start(ctx)
	if startErr != nil {
		WriteStartMetric(a.mainComponent, metrics.StartStatusFailed, startTime, a.metricsProvider)

		err := fmt.Errorf("%w: %v", ErrComponentStarting, startErr.Error())
		a.logger.Error("Start app components error", err)

		return err
	}

	WriteStartMetric(a.mainComponent, metrics.StartStatusSuccess, startTime, a.metricsProvider)

	return nil
}

//nolint:gochecknoglobals
var (
	scrapeMetricsTimeout = 20 * time.Second       // What is the maximum waiting time for sending metrics?
	scrapeMetricsRetry   = 500 * time.Millisecond // How often should I check what metrics were scrapped?
)

// Shutdown terminates the application and stops all functionality, starting with the main component.
func (a *App) Shutdown(ctx context.Context) error {
	var errs []error

	stopTime := time.Now().UTC()
	status := metrics.StopStatusSuccess
	before := atomic.LoadUint64(&metrics.LastScrapeCount)

	// In this implementation, this block is redundant because IComponent always contains the Shutdown function.
	// We could make IComponent only IStopper and leave IShutdowner as a separate interface,
	// but that idea seems strange.
	shtService, ok := a.mainComponent.(IShutdowner)
	if !ok {
		stopErr := a.mainComponent.Stop()
		if stopErr != nil {
			status = metrics.StopStatusFailed

			errs = append(errs, fmt.Errorf("%w: %v", ErrComponentStoping, stopErr.Error()))
		}
	}

	shutdownErr := shtService.Shutdown(ctx)
	if shutdownErr != nil {
		status = metrics.StopStatusNonSuccess

		err := fmt.Errorf("%w: %v", ErrComponentShutdowning, shutdownErr.Error())

		// if errors.Is(err, context.DeadlineExceeded) {
		// Special error, the component did not manage to stop within the allotted time.
		//
		// It might also be possible to add an error with "soft stop is not implemented"
		// in Shutdown when the component only has a Stop implementation.
		// }

		stopErr := a.mainComponent.Stop()
		if stopErr != nil {
			status = metrics.StopStatusFailed

			err = errors.Join(
				fmt.Errorf("%w: %v", ErrComponentShutdowning, shutdownErr.Error()),
				fmt.Errorf("%w: %v", ErrComponentStoping, stopErr.Error()),
			)
		}

		errs = append(errs, err)
	}

	WriteStopMetric(a.mainComponent, status, stopTime, a.metricsProvider)

	if len(errs) > 0 {
		a.logger.Error("Shutdown app components error", errors.Join(errs...))
	}

	// wait for scrape metrics
	deadline := time.Now().Add(scrapeMetricsTimeout)
	for time.Now().Before(deadline) {
		if atomic.LoadUint64(&metrics.LastScrapeCount) > before {
			break
		}

		time.Sleep(scrapeMetricsRetry)
	}

	// stopping the diag server is tied to the ctx of the main application,
	// but it is necessary to do it on the basis of a separate ctx.
	shutdownErr = a.diagServer.Shutdown(ctx)
	if shutdownErr != nil {
		// if errors.Is(err, context.DeadlineExceeded) {
		// Special error, the component did not manage to stop within the allotted time.
		//
		// It might also be possible to add an error with "soft stop is not implemented"
		// in Shutdown when the component only has a Stop implementation.
		// }
		closeErr := a.diagServer.Close()
		if closeErr != nil {
			a.logger.Error("Close server error", closeErr)
		}
	}

	return nil
}
