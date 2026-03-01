package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform"
	slicehelpers "github.com/mr-filatik/go-password-keeper/internal/platform/helpers/slice"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

type SequentialComponent struct { // Serial m.b.
	services []IComponent
	logger   logging.Logger

	// isStopped bool // для понимания, всё ли остановлено
}

// WithSequentialComponent combines components to start and stop them sequentially.
func (a *App) WithSequentialComponent(services ...IComponent) *SequentialComponent {
	if len(services) == 0 {
		services = append(services, &nopComponent{})
	}

	return &SequentialComponent{
		services: services,
		logger:   a.logger,
	}
}

// GetName displays the name of the component.
func (s *SequentialComponent) GetName() string {
	return "sequential component"
}

// Start starts the component.
func (s *SequentialComponent) Start(ctx context.Context) error {
	for _, service := range s.services {
		startErr := service.Start(ctx)
		if startErr != nil {
			// возвращается на первой ошибке, норма ли это?
			return fmt.Errorf("error when starting services sequentially: %w", startErr)
		}
	}

	return nil
}

// Shutdown initiates a soft stop of the component.
func (s *SequentialComponent) Shutdown(ctx context.Context) error {
	var errs []error

	reversed := slicehelpers.Reverse(s.services)

	for _, service := range reversed {
		shtService, ok := service.(platform.IShutdowner)
		if !ok {
			stopErr := service.Stop()
			if stopErr != nil {
				// возвращается на первой ошибке, норма ли это?
				return fmt.Errorf("error when stoping services sequentially: %w", stopErr)
			}
		}

		shutdownErr := shtService.Shutdown(ctx)
		if shutdownErr != nil {
			if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
				//s.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
			} else {
				//s.logger.Error("Shutdown component error", shutdownErr)
			}

			stopErr := service.Stop()
			if stopErr != nil {
				// возвращается на первой ошибке, норма ли это?
				return fmt.Errorf("error when stoping services sequentially: %w", stopErr)
			}

			return nil
		}
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
