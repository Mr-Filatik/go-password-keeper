// Package logctx provides utilities for embedding and retrieving logger instances
// within a context.Context using a strongly-typed key.
package logctx

import (
	"context"

	commonctx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

//nolint:gochecknoglobals
var ctxLoggerKey = &commonctx.CtxKey{Name: "logger"}

// SetLogger returns a new context containing the provided logger instance.
func SetLogger(ctx context.Context, value log.ILogger) context.Context {
	return commonctx.SetValue(ctx, ctxLoggerKey, value)
}

// GetLogger extracts the logger instance from the context.
// It returns a no-op logger (log.Nop) if the logger is not found or is nil.
//
//nolint:ireturn
func GetLogger(ctx context.Context) log.ILogger {
	logger, ok := commonctx.GetValue[log.ILogger](ctx, ctxLoggerKey)

	if !ok || logger == nil {
		return log.Nop
	}

	return logger
}
