// Package platform contains common functionality for all projects.
package platform

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// IStarter is an interface for all startup components.
type IStarter interface {
	// Start starts the component.
	Start(ctx context.Context) error
}

// Starter is the main entity for launching application components.
type Starter struct {
	logger   logging.Logger
	service  IStarter
	metrProv *metrics.Provider
}

// NewStarter creates a new instance of *Starter.
//
// Parameters:
//   - logger logging.Logger: logger;
//   - service IStarter: service for launching;
//   - metrProv *metrics.Provider: provider for tracking application launch metrics.
//
// The specified service is the main component from which the launch begins.
//
// To launch components in a specific order, it is recommended to use:
//   - SequentialStarter - to launch a list of services sequentially;
//   - ParallelStarter - to launch a list of services in parallel.
func NewStarter(logger logging.Logger, service IStarter, metrProv *metrics.Provider) *Starter {
	return &Starter{
		logger:   logger,
		service:  service,
		metrProv: metrProv,
	}
}

// Start starts the components.
//
// Parameters:
//   - ctx context.Context: context.
//
// Implements the platform.IStarter interface.
func (s *Starter) Start(ctx context.Context) error {
	startTime := time.Now().UTC()

	serviceStartErr := s.service.Start(ctx)
	if serviceStartErr != nil {
		s.metrProv.App.SetStartDuration(metrics.AppStartLabel{
			Status: metrics.StartStatusFailed,
		}, time.Since(startTime))

		s.logger.Error("Starting service error", serviceStartErr)

		return fmt.Errorf("error starting services: %w", serviceStartErr)
	}

	s.metrProv.App.SetStartDuration(metrics.AppStartLabel{
		Status: metrics.StartStatusSuccess,
	}, time.Since(startTime))

	// s.logger.Info("Starting service is successful")

	return nil
}

// SequentialStarter is an entity for sequentially starting application components.
type SequentialStarter struct {
	services []IStarter
}

// NewSequentialStarter creates a new *SequentialStarter instance for sequential starting.
//
// Parameters:
//   - services ...IStarter: list of services.
func NewSequentialStarter(services ...IStarter) *SequentialStarter {
	return &SequentialStarter{
		services: services,
	}
}

func (s *SequentialStarter) GetServices() []IStarter {
	return s.services
}

// Start initiates a sequential start of components.
//
// Parameters:
//   - ctx context.Context: context.
//
// Implements the platform.IStarter interface.
func (s *SequentialStarter) Start(ctx context.Context) error {
	for _, service := range s.services {
		startErr := service.Start(ctx)
		if startErr != nil {
			return fmt.Errorf("error when starting services sequentially: %w", startErr)
		}
	}

	return nil
}

// ParallelStarter is an entity for parallel launch of application components.
type ParallelStarter struct {
	services []IStarter
}

// NewParallelStarter creates a new *ParallelStarter instance to run in parallel.
//
// Parameters:
//   - services ...IStarter: list of services.
func NewParallelStarter(services ...IStarter) *ParallelStarter {
	return &ParallelStarter{
		services: services,
	}
}

func (s *ParallelStarter) GetServices() []IStarter {
	return s.services
}

// Start initiates parallel start of components.
//
// Parameters:
//   - ctx context.Context: context.
//
// Implements the platform.IStarter interface.
func (s *ParallelStarter) Start(ctx context.Context) error {
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
