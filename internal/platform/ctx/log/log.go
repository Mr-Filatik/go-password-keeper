package logctx

import (
	"context"

	tracectx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/trace"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// Debug logs a message at debug level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Debug(ctx context.Context, msg string, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Debug(msg, options...)
}

// Info logs a message at info level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Info(ctx context.Context, msg string, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Info(msg, options...)
}

// Warn logs a message and an error at warn level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Warn(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Warn(msg, err, options...)
}

// Error logs a message and an error at error level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Error(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Error(msg, err, options...)
}

// Fatal logs a message and an error at fatal level, then terminates the process.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Fatal(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Fatal(msg, err, options...)
}

//nolint:ireturn
func prepareLoggerAndOptions(ctx context.Context, options ...log.FieldOption) (log.ILogger, []log.FieldOption) {
	logger := GetLogger(ctx)

	autoLogger, ok := logger.(log.IAutoLogger)
	if ok && autoLogger.IsAutoTracing() {
		trace, ok := tracectx.GetTrace(ctx)
		if ok {
			options = append(options, log.WithAdvancedTraceField(trace)...)
		}
	}

	return logger, options
}
