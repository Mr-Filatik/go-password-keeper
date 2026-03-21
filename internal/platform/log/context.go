// Package logging provides logging functionality.
package log

import (
	"context"

	contextpkg "github.com/mr-filatik/go-password-keeper/internal/platform/context"
)

// ctxLoggerKey structure for storing and searching for a logger in a context.
//
//nolint:gochecknoglobals
var ctxLoggerKey = &contextpkg.CtxKey{Name: "logger"}

// ToContext sets the logger to the context.
func ToContext(ctx context.Context, logger ILogger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if logger == nil {
		logger = nop
	}

	return context.WithValue(ctx, ctxLoggerKey, logger)
}

// FromContext gets the logger from the context.
//
//nolint:ireturn // Necessary to fix the function in the interface.
func FromContext(ctx context.Context) ILogger {
	if ctx == nil {
		return nop
	}

	value := ctx.Value(ctxLoggerKey)
	if value == nil {
		return nop
	}

	logger, ok := value.(ILogger)
	if !ok || logger == nil {
		return nop
	}

	return logger
}

// Можно сделать ещё метод такой, чтобы caller был везде верный.
// Тогда надо запретить вызывать логгер напрямую. Но как, приватные поля?
func LogDebug(logger ILogger, msg string, options ...FieldOption) {
	if logger == nil {
		logger = nop
	}

	logger.Debug(msg, options...)
}

// LogInfo writes a log to the logger located in the context with the info level.
func LogInfo(logger ILogger, msg string, options ...FieldOption) {
	if logger == nil {
		logger = nop
	}

	logger.Info(msg, options...)

}

// LogWarn writes a log to the logger located in the context with the warning level.
func LogWarn(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = nop
	}

	logger.Warn(msg, err, options...)
}

// LogError writes a log to the logger located in the context with the error level.
func LogError(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = nop
	}

	logger.Error(msg, err, options...)
}

// LogFatal writes a log to the logger located in the context with the fatal level.
func LogFatal(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = nop
	}

	logger.Fatal(msg, err, options...)
}

// CtxDebug writes a log to the logger located in the context with the debug level.
func CtxDebug(ctx context.Context, msg string, options ...FieldOption) {
	FromContext(ctx).Debug(msg, options...)
}

// CtxInfo writes a log to the logger located in the context with the info level.
func CtxInfo(ctx context.Context, msg string, options ...FieldOption) {
	FromContext(ctx).Info(msg, options...)
}

// CtxWarn writes a log to the logger located in the context with the warning level.
func CtxWarn(ctx context.Context, msg string, err error, options ...FieldOption) {
	FromContext(ctx).Warn(msg, err, options...)
}

// CtxError writes a log to the logger located in the context with the error level.
func CtxError(ctx context.Context, msg string, err error, options ...FieldOption) {
	FromContext(ctx).Error(msg, err, options...)
}

// CtxFatal writes a log to the logger located in the context with the fatal level.
func CtxFatal(ctx context.Context, msg string, err error, options ...FieldOption) {
	FromContext(ctx).Fatal(msg, err, options...)
}
