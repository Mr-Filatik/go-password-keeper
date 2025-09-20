// Package logging provides logging functionality.
package logging

// LogFormat — формат вывода логов.
type LogFormat string

// Constants - описание типа вывода логов.
const (
	// FormatJSON - вывод логов в формате JSON.
	FormatJSON LogFormat = "JSON"

	// FormatText - вывод логов в текстовом формате.
	FormatText LogFormat = "TEXT"
)
