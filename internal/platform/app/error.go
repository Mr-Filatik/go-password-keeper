package app

import "errors"

// Common errors that can occur when starting or stopping application components.
var (
	ErrComponentStarting    = errors.New("error starting components")
	ErrComponentShutdowning = errors.New("error shutdowning components")
	ErrComponentStoping     = errors.New("error stoping components")
)
