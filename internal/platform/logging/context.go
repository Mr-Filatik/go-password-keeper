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

// Можно сделать ещё метод такой, чтобы caller был везде верный.
// Тогда надо запретить вызывать логгер напрямую. Но как, приватные поля?
func LogDebug(logger Logger, msg string, options ...FieldOption) {
	if logger != nil {
		logger.Debug(msg, options...)
	}
}

// LogInfo writes a log to the logger located in the context with the info level.
func LogInfo(logger Logger, msg string, options ...FieldOption) {
	if logger != nil {
		logger.Info(msg, options...)
	}
}

// LogWarn writes a log to the logger located in the context with the warning level.
func LogWarn(logger Logger, msg string, err error, options ...FieldOption) {
	if logger != nil {
		logger.Warn(msg, err, options...)
	}
}

// LogError writes a log to the logger located in the context with the error level.
func LogError(logger Logger, msg string, err error, options ...FieldOption) {
	if logger != nil {
		logger.Error(msg, err, options...)
	}
}

// LogFatal writes a log to the logger located in the context with the fatal level.
func LogFatal(logger Logger, msg string, err error, options ...FieldOption) {
	if logger != nil {
		logger.Fatal(msg, err, options...)
	}
}

// CtxDebug writes a log to the logger located in the context with the debug level.
func CtxDebug(ctx context.Context, msg string, options ...FieldOption) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Debug(msg, options...)
	}
}

// CtxInfo writes a log to the logger located in the context with the info level.
func CtxInfo(ctx context.Context, msg string, options ...FieldOption) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Info(msg, options...)
	}
}

// CtxWarn writes a log to the logger located in the context with the warning level.
func CtxWarn(ctx context.Context, msg string, err error, options ...FieldOption) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Warn(msg, err, options...)
	}
}

// CtxError writes a log to the logger located in the context with the error level.
func CtxError(ctx context.Context, msg string, err error, options ...FieldOption) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Error(msg, err, options...)
	}
}

// CtxFatal writes a log to the logger located in the context with the fatal level.
func CtxFatal(ctx context.Context, msg string, err error, options ...FieldOption) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Fatal(msg, err, options...)
	}
}
