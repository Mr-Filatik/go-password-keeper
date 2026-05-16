package log

import "io"

// ILogger describes the interface for all loggers used in the project.
//
// It is an implementation of the adapter pattern for converting any logger to a common interface.
type ILogger interface {
	// With returns a new logger with added common fields (labels).
	//
	// Parameters:
	//   - keysAndValues ...any: common fields (labels).
	//
	// Example:
	//   log := baseLogger.With("service", "issuetickets", "env", "test07")
	//   log.Info("something happened", "order_id", orderID)
	With(options ...FieldOption) ILogger

	// Debug logs the message and parameters with the debug level.
	//
	// Parameters:
	//   - message: main log message;
	//   - datas: additional information as a key-value pair.
	Debug(message string, options ...FieldOption)

	// Info logs the message and parameters with the info level.
	//
	// Parameters:
	//   - message: main log message;
	//   - datas: additional information as a key-value pair.
	Info(message string, options ...FieldOption)

	// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: possible error;
	//   - datas: additional information as a key-value pair.
	Warn(message string, err error, options ...FieldOption)

	// Error logs a message and parameters with the error level and error.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: error;
	//   - datas: additional information as a key-value pair.
	Error(message string, err error, options ...FieldOption)

	// Fatal logs a message and parameters with the fatal and critical error levels.
	//
	// Parameters:
	//   - message: main log message;
	//   - err: critical error;
	//   - datas: additional information as a key-value pair.
	Fatal(message string, err error, options ...FieldOption)

	// Close releases resources used by the logger.
	//
	// Implements the io.Closer interface.
	io.Closer
}

type IAutoLogger interface {
	IsAutoTracing() bool
}
