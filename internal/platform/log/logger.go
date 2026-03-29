// Package log provides logging functionality.
package log

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

// OutputFormat — log output format.
type OutputFormat string

// Constants - description of the log output type.
const (
	// FormatJSON - output logs in JSON format.
	OutputFormatJSON OutputFormat = "JSON"

	// FormatText - output logs in text format.
	OutputFormatText OutputFormat = "TEXT"
)

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
