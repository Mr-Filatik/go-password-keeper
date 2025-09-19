// Package logging provides logging functionality.
package logging

import "io"

// LogLevel - logging level.
type LogLevel uint8

// Constants - specific logging levels.
const (
	// LevelDebug - debug logging level.
	LevelDebug LogLevel = iota

	// LevelInfo - info logging level.
	LevelInfo

	// LevelWarn - warning logging level.
	LevelWarn

	// LevelError - error logging level.
	LevelError

	// LevelFatal - fatal logging level.
	LevelFatal
)

// Logger describes the interface for all loggers used in the project.
//
// It is an implementation of the adapter pattern for converting any logger to a common interface.
type Logger interface {
	// Debug logs the message and parameters with the debug level.
	Debug(message string, keysAndValues ...any)

	// Info logs the message and parameters with the info level.
	Info(message string, keysAndValues ...any)

	// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
	Warn(message string, err error, keysAndValues ...any)

	// Error logs a message and parameters with the error level and error.
	Error(message string, err error, keysAndValues ...any)

	// Fatal logs a message and parameters with the fatal and critical error levels.
	Fatal(message string, err error, keysAndValues ...any)

	// Close releases resources used by the logger.
	//
	// Implements the io.Closer interface.
	io.Closer
}

// String outputs a string representation corresponding to the logging level.
//
// Implements the fmt.Stringer interface.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "none"
	}
}

// Limit limits the level to the maximum allowed if an invalid value is specified.
func (l LogLevel) Limit() LogLevel {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal:
		return l
	default:
		return LevelFatal
	}
}
