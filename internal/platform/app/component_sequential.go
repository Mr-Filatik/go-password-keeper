package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	slicehelpers "github.com/mr-filatik/go-password-keeper/internal/platform/helpers/slice"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// SequentialComponent represents a special component used to sequentially start child components.
// They are stopped in reverse order.
type SequentialComponent struct {
	components      []IComponent
	metricsProvider IAppMetrics

	stopLaunchingOnError bool // whether to stop the launch when the first error is received
}

// WithSequentialComponent creates a component for sequentially starting and stopping elements.
// If no components are passed, one nested nopComponent will be created.
//
// Created based on the App structure for sharing dependencies.
func (a *App) WithSequentialComponent(services ...IComponent) *SequentialComponent {
	if len(services) == 0 {
		services = append(services, &nopComponent{a.logger})
	}

	return &SequentialComponent{
		components:           services,
		metricsProvider:      a.metricsProvider,
		stopLaunchingOnError: a.stopLaunchingOnError,
	}
}

// GetName displays the name of the component.
//
// Implements the IComponent interface.
func (c *SequentialComponent) GetName() string {
	return "sequential component"
}

// Start begins launching all components within the SequentialComponent.
//
// If the app.WithStopLaunchingOnError() option was specified when creating the App,
// then the component launch will be interrupted by the first error encountered.
//
// Implements the IComponent interface.
func (c *SequentialComponent) Start(ctx context.Context) error {
	var errs []error

	for _, service := range c.components {
		startTime := time.Now().UTC()
		status := metrics.StartStatusSuccess

		startErr := service.Start(ctx)
		if startErr != nil {
			status = metrics.StartStatusFailed

			if c.stopLaunchingOnError {
				WriteStartMetric(service, status, startTime, c.metricsProvider)

				return fmt.Errorf("starting sequential components: %w", startErr)
			}

			errs = append(errs, startErr)
		}

		WriteStartMetric(service, status, startTime, c.metricsProvider)
	}

	if len(errs) > 0 {
		return fmt.Errorf("starting sequential components: %w", errors.Join(errs...))
	}

	return nil
}

// Shutdown function initiates a soft stop of all components in the SequentialComponent.
//
// If an error occurs when calling Shutdown for an internal component
// (or if it doesn't stop within the allotted time), Stop is called.
//
// Implements the IComponent interface.
func (c *SequentialComponent) Shutdown(ctx context.Context) error {
	var errs []error

	reversed := slicehelpers.Reverse(c.components)

	for _, service := range reversed {
		stopTime := time.Now().UTC()
		status := metrics.StopStatusSuccess

		// In this implementation, this block is redundant because IComponent always contains the Shutdown function.
		// We could make IComponent only IStopper and leave IShutdowner as a separate interface,
		// but that idea seems strange.
		shtService, ok := service.(IShutdowner)
		if !ok {
			stopErr := service.Stop()
			if stopErr != nil {
				status = metrics.StopStatusFailed

				errs = append(errs, stopErr)
			}
		}

		shutdownErr := shtService.Shutdown(ctx)
		if shutdownErr != nil {
			status = metrics.StopStatusNonSuccess

			err := shutdownErr

			// if errors.Is(err, context.DeadlineExceeded) {
			// Special error, the component did not manage to stop within the allotted time.
			//
			// It might also be possible to add an error with "soft stop is not implemented"
			// in Shutdown when the component only has a Stop implementation.
			// }

			stopErr := service.Stop()
			if stopErr != nil {
				status = metrics.StopStatusFailed

				err = errors.Join(shutdownErr, stopErr)
			}

			errs = append(errs, err)
		}

		WriteStopMetric(service, status, stopTime, c.metricsProvider)
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdowning sequential components: %w", errors.Join(errs...))
	}

	return nil
}

// Stop function initiates stopping all components in the SequentialComponent.
//
// Implements the IComponent interface.
func (c *SequentialComponent) Stop() error {
	var errs []error

	reversed := slicehelpers.Reverse(c.components)

	for _, service := range reversed {
		stopTime := time.Now().UTC()
		status := metrics.StopStatusSuccess

		stopErr := service.Stop()
		if stopErr != nil {
			status = metrics.StopStatusFailed

			errs = append(errs, stopErr)
		}

		WriteStopMetric(service, status, stopTime, c.metricsProvider)
	}

	if len(errs) > 0 {
		return fmt.Errorf("stopping sequential components: %w", errors.Join(errs...))
	}

	return nil
}
