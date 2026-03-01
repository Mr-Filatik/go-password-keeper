package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

type ParallelComponent struct { // Serial m.b.
	services []IComponent
	logger   logging.Logger
}

// WithParallelComponent combines components to start and stop in parallel.
func (a *App) WithParallelComponent(services ...IComponent) *ParallelComponent {
	if len(services) == 0 {
		services = append(services, &nopComponent{})
	}

	return &ParallelComponent{
		services: services,
		logger:   a.logger,
	}
}

// GetName displays the name of the component.
func (s *ParallelComponent) GetName() string {
	return "parallel component"
}

// Start starts the component.
func (s *ParallelComponent) Start(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.services))

	for _, service := range s.services {
		go func() {
			defer wg.Done()

			startErr := service.Start(ctx)
			if startErr != nil {
				mu.Lock()

				errs = append(errs, startErr)

				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("error when running services in parallel: %w", errors.Join(errs...))
	}

	return nil
}

// Shutdown initiates a soft stop of the component.
func (s *ParallelComponent) Shutdown(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.services))

	for _, service := range s.services {
		go func() {
			defer wg.Done()

			shtService, ok := service.(IShutdowner)
			if ok {
				shutdownErr := shtService.Shutdown(ctx)
				if shutdownErr != nil {
					if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
						s.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
					} else {
						s.logger.Error("Shutdown component error", shutdownErr)
					}

					stopErr := service.Stop()
					if stopErr != nil {
						s.logger.Error("Stoping service error", stopErr)

						mu.Lock()

						errs = append(errs, stopErr)

						mu.Unlock()
					}
				}
			} else {
				stopErr := service.Stop()
				if stopErr != nil {
					s.logger.Error("Stoping service error", stopErr)

					mu.Lock()

					errs = append(errs, stopErr)

					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	s.logger.Warn("Shutdowning is done", nil, "component", s.GetName()) // TEMP

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Stop starts stopping the component.
func (s *ParallelComponent) Stop() error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(s.services))

	for _, service := range s.services {
		go func() {
			defer wg.Done()

			stopErr := service.Stop()
			if stopErr != nil {
				s.logger.Error("Stoping service error", stopErr)

				mu.Lock()

				errs = append(errs, stopErr)

				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
