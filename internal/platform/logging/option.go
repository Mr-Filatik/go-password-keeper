package logging

import (
	"io"
	"os"
)

type config struct {
	// logLevel        LogLevel
	writer          io.Writer
	format          LogFormat
	callerSkipCount int

	globalFields []any
}

func defaultConfig() config {
	return config{
		writer:          os.Stdout,
		format:          FormatJSON,
		callerSkipCount: 2,

		globalFields: []any{},
	}
}

type ConfigOption func(config *config)

func WithWriter(writer io.Writer) ConfigOption {
	return func(config *config) {
		config.writer = writer
	}
}

func WithFormat(format LogFormat) ConfigOption {
	return func(config *config) {
		config.format = format
	}
}

// TODO zap.AddCallerSkip(1) если напрямую логгер
// TODO zap.AddCallerSkip(2) если логгер используется через Log... Ctx...
func WithCallerSkip(count int) ConfigOption {
	return func(config *config) {
		config.callerSkipCount = count
	}
}

func WithGlobalFields(options ...FieldOption) ConfigOption {
	return func(config *config) {
		config.globalFields = applyOptions(options...)
	}
}
