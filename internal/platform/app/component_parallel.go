package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// ParallelComponent represents a special component used to parallel start child components.
type ParallelComponent struct {
	components      []IComponent
	logger          logging.Logger
	metricsProvider IAppMetrics
}

// WithParallelComponent creates a component for parallel starting and stopping elements.
// If no components are passed, one nested nopComponent will be created.
//
// Created based on the App structure for sharing dependencies.
func (a *App) WithParallelComponent(services ...IComponent) *ParallelComponent {
	if len(services) == 0 {
		services = append(services, &nopComponent{a.logger})
	}

	return &ParallelComponent{
		components:      services,
		logger:          a.logger,
		metricsProvider: a.metricsProvider,
	}
}

// GetName displays the name of the component.
//
// Implements the IComponent interface.
func (s *ParallelComponent) GetName() string {
	return "parallel component"
}

// Start begins launching all components within the ParallelComponent.
//
// Implements the IComponent interface.
func (s *ParallelComponent) Start(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.components))

	for _, service := range s.components {
		go func() {
			defer wg.Done()

			startTime := time.Now().UTC()
			status := metrics.StartStatusSuccess

			startErr := service.Start(ctx)
			if startErr != nil {
				status = metrics.StartStatusFailed

				mu.Lock()

				errs = append(errs, startErr)

				mu.Unlock()
			}

			WriteStartMetric(service, status, startTime, s.metricsProvider)
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("starting parallel components: %w", errors.Join(errs...))
	}

	return nil
}

// Shutdown function initiates a soft stop of all components in the ParallelComponent.
//
// If an error occurs when calling Shutdown for an internal component
// (or if it doesn't stop within the allotted time), Stop is called.
//
// Implements the IComponent interface.
//
//nolint:funlen // Comments in the function are still needed
func (s *ParallelComponent) Shutdown(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.components))

	for _, service := range s.components {
		go func() {
			defer wg.Done()

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

					mu.Lock()

					errs = append(errs, stopErr)

					mu.Unlock()
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

				mu.Lock()

				errs = append(errs, err)

				mu.Unlock()
			}

			WriteStopMetric(service, status, stopTime, s.metricsProvider)
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("shutdowning parallel components: %w", errors.Join(errs...))
	}

	return nil
}

// Stop function initiates stopping all components in the ParallelComponent.
//
// Implements the IComponent interface.
func (s *ParallelComponent) Stop() error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.components))

	for _, service := range s.components {
		go func() {
			defer wg.Done()

			stopTime := time.Now().UTC()
			status := metrics.StopStatusSuccess

			stopErr := service.Stop()
			if stopErr != nil {
				status = metrics.StopStatusFailed

				mu.Lock()

				errs = append(errs, stopErr)

				mu.Unlock()
			}

			WriteStopMetric(service, status, stopTime, s.metricsProvider)
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("stopping parallel components: %w", errors.Join(errs...))
	}

	return nil
}
