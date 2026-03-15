// Package logging provides logging functionality.
package logging

import (
	"context"

	contextpkg "github.com/mr-filatik/go-password-keeper/internal/platform/context"
)

// CtxKeyLogger structure for storing and searching for a logger in a context.
//
//nolint:gochecknoglobals // Migrate from platform to http
var CtxKeyLogger = &contextpkg.CtxKey{Name: "logger"}

// ToContext sets the logger to the context.
func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, CtxKeyLogger, logger)
}

// FromContext gets the logger from the context.
//
//nolint:ireturn // Necessary to fix the function in the interface.
func FromContext(ctx context.Context) Logger {
	value := ctx.Value(CtxKeyLogger)
	if value == nil {
		return nil
	}

	logger, ok := value.(Logger)
	if !ok {
		return nil
	}

	return logger
}

func Debug(ctx context.Context, msg string, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Debug(msg, datas...)
	}
}

func Info(ctx context.Context, msg string, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Info(msg, datas...)
	}
}

func Warn(ctx context.Context, msg string, err error, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Warn(msg, err, datas...)
	}
}

func Error(ctx context.Context, msg string, err error, datas ...any) {
	logger := FromContext(ctx)

	if logger != nil {
		logger.Error(msg, err, datas...)
	}
}
