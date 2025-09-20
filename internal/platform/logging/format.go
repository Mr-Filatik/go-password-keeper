// Package logging provides logging functionality.
package logging

// LogFormat — log output format.
type LogFormat string

// Constants - description of the log output type.
const (
	// FormatJSON - output logs in JSON format.
	FormatJSON LogFormat = "JSON"

	// FormatText - output logs in text format.
	FormatText LogFormat = "TEXT"
)
