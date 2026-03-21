package log

import (
	"io"
	"os"
)

type Config struct {
	// logLevel        LogLevel
	writer          io.Writer
	format          LogFormat
	callerSkipCount int

	globalFields []any
}

func DefaultConfig() Config {
	return Config{
		writer:          os.Stdout,
		format:          FormatJSON,
		callerSkipCount: 2,

		globalFields: []any{},
	}
}

func (c Config) GetWriter() io.Writer {
	return c.writer
}

func (c Config) GetFormat() LogFormat {
	return c.format
}

func (c Config) GetCallerSkipCount() int {
	return c.callerSkipCount
}

func (c Config) GetGlobalFields() []any {
	if c.globalFields == nil {
		return nil
	}

	fields := make([]any, len(c.globalFields))
	copy(fields, c.globalFields)

	return c.globalFields
}

type ConfigOption func(config *Config)

func WithWriter(writer io.Writer) ConfigOption {
	return func(config *Config) {
		config.writer = writer
	}
}

func WithFormat(format LogFormat) ConfigOption {
	return func(config *Config) {
		config.format = format
	}
}

// TODO zap.AddCallerSkip(1) если напрямую логгер
// TODO zap.AddCallerSkip(2) если логгер используется через Log... Ctx...
func WithCallerSkip(count int) ConfigOption {
	return func(config *Config) {
		config.callerSkipCount = count
	}
}

func WithGlobalFields(options ...FieldOption) ConfigOption {
	return func(config *Config) {
		config.globalFields = ApplyOptions(options...)
	}
}
