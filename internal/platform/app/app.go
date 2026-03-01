// Package app provides everything needed to work with application components.
package app

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/diagnostic"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// App provides a framework for starting and stopping application containers.
type App struct {
	logger        logging.Logger
	mainComponent IComponent
	metrProv      *metrics.Provider
	diagServer    *diagnostic.Server
}

// New creates an instance of the App structure.
func New(logger logging.Logger, metrProv *metrics.Provider, diagServer *diagnostic.Server) *App {
	return &App{
		logger:        logger,
		mainComponent: &nopComponent{},
		metrProv:      metrProv,
		diagServer:    diagServer,
	}
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

// Start starts launching all application components, starting with the main one.
func (a *App) Start(ctx context.Context) error {
	startTime := time.Now().UTC()

	startErr := a.mainComponent.Start(ctx)
	if startErr != nil {
		a.metrProv.App.SetStartDuration(metrics.AppStartLabel{
			Status: metrics.StartStatusFailed,
		}, time.Since(startTime))

		return fmt.Errorf("%w: %v", ErrComponentStarting, startErr.Error())
	}

	// нужно сделать так, чтобы каждый отдельный компонент писался в метрику
	// но нужно исключать повторные запуски, компоненты паралел и сиквеншиал, само App
	// они дублируют основную информацию
	a.metrProv.App.SetStartDuration(metrics.AppStartLabel{
		Status: metrics.StartStatusSuccess,
	}, time.Since(startTime))

	return nil
}

// Shutdown softly stops all application components, in reverse order.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Debug("App shutdowning...")

	stopTime := time.Now().UTC()
	before := atomic.LoadUint64(&metrics.LastScrapeCount)

	// var (
	// 	wg   sync.WaitGroup
	// 	mu   sync.Mutex
	// 	errs []error
	// )

	shtService, ok := a.mainComponent.(IShutdowner)
	if !ok {
		stopErr := a.mainComponent.Stop()
		if stopErr != nil {
			a.metrProv.App.SetStopDuration(metrics.AppStopLabel{
				Status: metrics.StopStatusFailed,
			}, time.Since(stopTime))

			return fmt.Errorf("%w: %v", ErrComponentStoping, stopErr.Error())
		}
	}

	shutdownErr := shtService.Shutdown(ctx)
	if shutdownErr != nil {
		if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
			a.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
		} else {
			a.logger.Error("Shutdown component error", shutdownErr)
		}

		stopErr := a.mainComponent.Stop()
		if stopErr != nil {
			a.metrProv.App.SetStopDuration(metrics.AppStopLabel{
				Status: metrics.StopStatusFailed,
			}, time.Since(stopTime))

			return fmt.Errorf("%w: %v", ErrComponentShutdowning, stopErr.Error())
		}
	}

	a.metrProv.App.SetStopDuration(metrics.AppStopLabel{
		Status: metrics.StopStatusSuccess,
	}, time.Since(stopTime))

	// logger.Info("Application shutdown is successful") // время

	// wait for scrape metrics
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadUint64(&metrics.LastScrapeCount) > before {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	shutdownErr = a.diagServer.Shutdown(ctx)
	if shutdownErr != nil {
		if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
			a.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
		} else {
			a.logger.Error("Shutdown component error", shutdownErr)
		}

		closeErr := a.diagServer.Close()
		if closeErr != nil {
			a.logger.Error("Close server error", closeErr)
		}
	}

	return nil
}

// func (a *App) Stop() error {
// 	a.logger.Debug("App stoping...")

// 	// logic

// 	return a.mainComponent.Stop()
// }
