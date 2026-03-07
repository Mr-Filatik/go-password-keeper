package app

import (
	"context"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// IComponent describes a universal interface for any application component.
// This could be a server, client, adapter, or other component.
// All these components must implement startup and shutdown (soft and hard).
type IComponent interface {
	// GetName displays the name of the component.
	GetName() string

	IStarter
	IShutdowner
}

// IStarter provides a universal interface for all startup components.
type IStarter interface {
	// Start starts the component.
	//
	// A bool type variable must be used inside the component
	// to prevent the component from being launched twice.
	Start(ctx context.Context) error
}

// IShutdowner provides a universal interface for all soft-stop components.
type IShutdowner interface {
	// Shutdown initiates a soft stop of the component.
	Shutdown(ctx context.Context) error

	IStopper
}

// IStopper provides a universal interface for all stoppable components.
type IStopper interface {
	// Stop starts stopping the component.
	Stop() error
}

// IAppMetrics describes an interface for metrics about component start and stop.
type IAppMetrics interface {
	SetStartDuration(labels metrics.AppStartLabel, duration time.Duration)
	SetStopDuration(labels metrics.AppStopLabel, duration time.Duration)
}
