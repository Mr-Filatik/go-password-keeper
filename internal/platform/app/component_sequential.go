package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform"
	slicehelpers "github.com/mr-filatik/go-password-keeper/internal/platform/helpers/slice"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

type SequentialComponent struct { // Serial m.b.
	services []IComponent
	logger   logging.Logger
	metrProv *metrics.Provider

	stopLaunchingOnError bool
	// isStopped bool // для понимания, всё ли остановлено
}

// WithSequentialComponent combines components to start and stop them sequentially.
func (a *App) WithSequentialComponent(services ...IComponent) *SequentialComponent {
	if len(services) == 0 {
		services = append(services, &nopComponent{a.logger})
	}

	return &SequentialComponent{
		services:             services,
		logger:               a.logger,
		metrProv:             a.metrProv,
		stopLaunchingOnError: a.stopLaunchingOnError,
	}
}

// GetName displays the name of the component.
//
// Implements the IComponent interface.
func (s *SequentialComponent) GetName() string {
	return "sequential component"
}

// Start begins launching all components within the SequentialComponent.
//
// If the app.WithStopLaunchingOnError() option was specified when creating the App,
// then the component launch will be interrupted by the first error encountered.
//
// Implements the IComponent interface.
func (s *SequentialComponent) Start(ctx context.Context) error {
	var errs []error

	for _, service := range s.services {
		startTime := time.Now().UTC()
		status := metrics.StartStatusSuccess

		startErr := service.Start(ctx)
		if startErr != nil {
			status = metrics.StartStatusFailed

			err := fmt.Errorf("error when starting services sequentially: %w", startErr)

			if s.stopLaunchingOnError {
				WriteStartMetric(service, metrics.StartStatusFailed, startTime, s.metrProv.App)

				return err
			}

			errs = append(errs, err)
		}

		WriteStartMetric(service, status, startTime, s.metrProv.App)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Shutdown initiates a soft stop of the component.
func (s *SequentialComponent) Shutdown(ctx context.Context) error {
	var errs []error

	reversed := slicehelpers.Reverse(s.services)

	for _, service := range reversed {
		stopTime := time.Now().UTC()
		status := metrics.StopStatusSuccess

		// этот вариант лишний, т.к. IComponent всегда содержит метод Shutdown
		// но можно удалить из IComponent Shutdown и оставить IStop, тогда Shutdown будет необязательным
		shtService, ok := service.(platform.IShutdowner)
		if !ok {
			stopErr := service.Stop()
			if stopErr != nil {
				status = metrics.StopStatusFailed

				errs = append(errs, fmt.Errorf("error when stoping services sequentially: %w", stopErr))
			}
		}

		shutdownErr := shtService.Shutdown(ctx)
		if shutdownErr != nil {
			status = metrics.StopStatusNonSuccess

			if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
				//s.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
			} else {
				//s.logger.Error("Shutdown component error", shutdownErr)
			}

			stopErr := service.Stop()
			if stopErr != nil {
				status = metrics.StopStatusFailed

				errs = append(errs, fmt.Errorf("error when stoping services sequentially: %w", stopErr))
			}
		}

		WriteStopMetric(service, status, stopTime, s.metrProv.App)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Stop starts stopping the component.
func (s *SequentialComponent) Stop() error {
	var errs []error

	reversed := slicehelpers.Reverse(s.services)

	for _, service := range reversed {
		stopErr := service.Stop()
		if stopErr != nil {
			// возвращается на первой ошибке, норма ли это?
			return fmt.Errorf("error when stoping services sequentially: %w", stopErr)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
