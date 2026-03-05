package app

import (
	"context"
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// nopComponent represents a stub component.
type nopComponent struct {
	logger logging.Logger
}

// GetName displays the name of the component.
func (c *nopComponent) GetName() string {
	return "nop component"
}

// Start starts the component.
func (c *nopComponent) Start(_ context.Context) error {
	err := fmt.Errorf("%w: nopComponent was launched", ErrComponentStarting)

	c.logger.Warn("Actions with a service object", err)

	return nil
}

// Shutdown initiates a soft stop of the component.
func (c *nopComponent) Shutdown(_ context.Context) error {
	err := fmt.Errorf("%w: nopComponent was soft stopped", ErrComponentShutdowning)

	c.logger.Warn("Actions with a service object", err)

	return nil
}

// Stop starts stopping the component.
func (c *nopComponent) Stop() error {
	err := fmt.Errorf("%w: nopComponent has been stopped", ErrComponentStoping)

	c.logger.Warn("Actions with a service object", err)

	return nil
}
