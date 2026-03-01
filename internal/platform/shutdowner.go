// Package platform contains common functionality for all projects.
package platform

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/http/diagnostic"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// replace io.Closer for IStopper
type IStopper interface {
	Stop() error
}

// IShutdowner is an interface for all shutdown components.
type IShutdowner interface {
	// Shutdown softly stops the component.
	Shutdown(ctx context.Context) error

	// // Close hard stops the component.
	// //
	// // Implements the io.Closer interface.
	// io.Closer

	IStopper
}

type Stopper struct {
	logger      logging.Logger
	service     IStopper
	metrProv    *metrics.Provider
	diagSerrver *diagnostic.Server
}

func NewStopper(
	logger logging.Logger,
	service IStopper,
	metrProv *metrics.Provider,
	diagSerrver *diagnostic.Server,
) *Stopper {
	return &Stopper{
		logger:      logger,
		service:     service,
		metrProv:    metrProv,
		diagSerrver: diagSerrver,
	}
}

func (s *Stopper) Shutdown(ctx context.Context) error {
	stopTime := time.Now().UTC()
	before := atomic.LoadUint64(&metrics.LastScrapeCount)

	// var (
	// 	wg   sync.WaitGroup
	// 	mu   sync.Mutex
	// 	errs []error
	// )

	shtService, ok := s.service.(IShutdowner)
	if ok {
		shutdownErr := shtService.Shutdown(ctx)
		if shutdownErr != nil {
			if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
				s.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
			} else {
				s.logger.Error("Shutdown component error", shutdownErr)
			}

			stopErr := s.service.Stop()
			if stopErr != nil {
				s.metrProv.App.SetStopDuration(metrics.AppStopLabel{
					Status: metrics.StopStatusFailed,
				}, time.Since(stopTime))

				s.logger.Error("Stoping service error", stopErr)

				return stopErr
			}
		}
	} else {
		stopErr := s.service.Stop()
		if stopErr != nil {
			s.metrProv.App.SetStopDuration(metrics.AppStopLabel{
				Status: metrics.StopStatusFailed,
			}, time.Since(stopTime))

			s.logger.Error("Stoping service error", stopErr)

			return stopErr
		}
	}

	s.metrProv.App.SetStopDuration(metrics.AppStopLabel{
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

	shutdownErr := s.diagSerrver.Shutdown(ctx)
	if shutdownErr != nil {
		if errors.Is(shutdownErr, context.DeadlineExceeded) { // ErrServerClosed ??
			s.logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
		} else {
			s.logger.Error("Shutdown component error", shutdownErr)
		}

		closeErr := s.diagSerrver.Close()
		if closeErr != nil {
			s.logger.Error("Close server error", closeErr)
		}
	}

	return nil
}

type SequentialStopper struct { // Serial m.b.
	logger   logging.Logger
	services []IStopper
}

func NewSequentialStopper(logger logging.Logger, services ...IStopper) *SequentialStopper {
	return &SequentialStopper{
		logger:   logger,
		services: services,
	}
}

func (s *SequentialStopper) Shutdown(ctx context.Context) error {
	var errs []error

	for _, service := range s.services {
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

					return stopErr
				}
			}
		} else {
			stopErr := service.Stop()
			if stopErr != nil {
				s.logger.Error("Stoping service error", stopErr)

				return stopErr
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *SequentialStopper) Stop() error {
	var errs []error

	for _, service := range s.services {
		stopErr := service.Stop()
		if stopErr != nil {
			s.logger.Error("Stoping service error", stopErr)

			return stopErr
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

type ParallelStopper struct {
	logger   logging.Logger
	services []IStopper
}

func NewParallelStopper(logger logging.Logger, services ...IStopper) *ParallelStopper {
	return &ParallelStopper{
		logger:   logger,
		services: services,
	}
}

func (s *ParallelStopper) Shutdown(ctx context.Context) error {
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

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *ParallelStopper) Stop() error {
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
