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

	// DefaultLevel - minimum log limit level.
	DefaultLevel = LevelError
)

// LogFormat — log output format.
type LogFormat string

// Constants - description of the log output type.
const (
	// FormatJSON - output logs in JSON format.
	FormatJSON LogFormat = "JSON"

	// FormatText - output logs in text format.
	FormatText LogFormat = "TEXT"
)

// Logger describes the interface for all loggers used in the project.
//
// It is an implementation of the adapter pattern for converting any logger to a common interface.
type Logger interface {
	// With returns a new logger with added common fields (labels).
	//
	// Parameters:
	//   - keysAndValues ...any: common fields (labels).
	//
	// Example:
	//   log := baseLogger.With("service", "issuetickets", "env", "test07")
	//   log.Info("something happened", "order_id", orderID)
	With(keysAndValues ...any) Logger

	// Debug logs the message and parameters with the debug level.
	//
	// Parameters:
	//   - message: main log message;
	//   - datas: additional information as a key-value pair.
	Debug(message string, datas ...any)

	// Info logs the message and parameters with the info level.
	//
	// Parameters:
	//   - message: main log message;
	//   - datas: additional information as a key-value pair.
	Info(message string, datas ...any)

	// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: possible error;
	//   - datas: additional information as a key-value pair.
	Warn(message string, err error, datas ...any)

	// Error logs a message and parameters with the error level and error.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: error;
	//   - datas: additional information as a key-value pair.
	Error(message string, err error, datas ...any)

	// Fatal logs a message and parameters with the fatal and critical error levels.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: critical error;
	//   - datas: additional information as a key-value pair.
	Fatal(message string, err error, datas ...any)

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

// Validate limits the level to the maximum allowed if an invalid value is specified.
func (l LogLevel) Validate() LogLevel {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return l
	case LevelFatal:
		return DefaultLevel
	default:
		return DefaultLevel
	}
}
