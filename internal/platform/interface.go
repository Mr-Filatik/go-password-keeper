// Package platform contains common functionality for all projects.
package platform

import (
	"context"
	"io"
)

// IStarter is an interface for all startup components.
type IStarter interface {
	// Start starts the component.
	Start(ctx context.Context) error
}

// IShutdowner is an interface for all shutdown components.
type IShutdowner interface {
	// Shutdown softly stops the component.
	Shutdown(ctx context.Context) error

	// Close hard stops the component.
	//
	// Implements the io.Closer interface.
	io.Closer
}
