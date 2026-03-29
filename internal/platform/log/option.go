package log

import (
	"io"
	"os"
)

type Config struct {
	// logLevel        LogLevel
	writer          io.Writer
	format          OutputFormat
	callerSkipCount int

	globalFields []any
}

func DefaultConfig() Config {
	return Config{
		writer:          os.Stdout,
		format:          OutputFormatJSON,
		callerSkipCount: int(CallerSkipIndirect),

		globalFields: []any{},
	}
}

func (c Config) GetWriter() io.Writer {
	return c.writer
}

func (c Config) GetFormat() OutputFormat {
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

func WithFormat(format OutputFormat) ConfigOption {
	return func(config *Config) {
		config.format = format
	}
}

type CallerSkip int

const (
	// CallerSkipNone если без логера вызываются
	CallerSkipNone CallerSkip = 0

	// CallerSkipDirect если напрямую логгер
	CallerSkipDirect CallerSkip = 1

	// CallerSkipIndirect если логгер используется через Log... Ctx...
	CallerSkipIndirect CallerSkip = 2
)

func WithCallerSkip(count CallerSkip) ConfigOption {
	return func(config *Config) {
		config.callerSkipCount = int(count)
	}
}

func WithGlobalFields(options ...FieldOption) ConfigOption {
	return func(config *Config) {
		config.globalFields = ApplyOptions(options...)
	}
}
