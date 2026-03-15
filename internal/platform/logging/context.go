// Package logging provides logging functionality.
package logging

import (
	"context"

	contextpkg "github.com/mr-filatik/go-password-keeper/internal/platform/context"
)

// ctxLoggerKey structure for storing and searching for a logger in a context.
//
//nolint:gochecknoglobals
var ctxLoggerKey = &contextpkg.CtxKey{Name: "logger"}

// ToContext sets the logger to the context.
func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

// FromContext gets the logger from the context.
//
//nolint:ireturn // Necessary to fix the function in the interface.
func FromContext(ctx context.Context) Logger {
	value := ctx.Value(ctxLoggerKey)
	if value == nil {
		return nil // return &noopLogger{}
	}

	logger, ok := value.(Logger)
	if !ok {
		return nil
	}

	return logger
}

// Debug writes a log to the logger located in the context with the debug level.
func Debug(ctx context.Context, msg string, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Debug(msg, datas...)
	}
}

// Info writes a log to the logger located in the context with the info level.
func Info(ctx context.Context, msg string, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Info(msg, datas...)
	}
}

// Warn writes a log to the logger located in the context with the warning level.
func Warn(ctx context.Context, msg string, err error, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Warn(msg, err, datas...)
	}
}

// Error writes a log to the logger located in the context with the error level.
func Error(ctx context.Context, msg string, err error, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Error(msg, err, datas...)
	}
}

// Fatal writes a log to the logger located in the context with the fatal level.
func Fatal(ctx context.Context, msg string, err error, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Fatal(msg, err, datas...)
	}
}
