package app

import (
	"context"
)

// nopComponent represents a stub component.
type nopComponent struct{}

// GetName displays the name of the component.
func (s *nopComponent) GetName() string {
	return "nop component"
}

// Start starts the component.
func (s *nopComponent) Start(_ context.Context) error {
	return nil
}

// Shutdown initiates a soft stop of the component.
func (s *nopComponent) Shutdown(_ context.Context) error {
	return nil
}

// Stop starts stopping the component.
func (s *nopComponent) Stop() error {
	return nil
}
