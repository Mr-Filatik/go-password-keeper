// Package context contains functionality for working with context.
package context

import (
	"context"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// CtxKey is a structure describing keys for working with context variables.
type CtxKey struct {
	// Name - name of the key.
	Name string
}

// WithValue sets the value (string) by key in the context.
//
// Parameters:
//   - ctx context.Context: parent context;
//   - key *CtxKey: key;
//   - value string: value as a string.
func WithValue(ctx context.Context, key *CtxKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue gets the value (string) from the context by key.
//
// Parameters:
// - ctx context.Context: context;
// - key *CtxKey: key.
func GetValue(ctx context.Context, key *CtxKey) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}

// CtxKeyLogger - key for the ctx header "logger".
//
//nolint:gochecknoglobals // Migrate from platform to http
var CtxKeyLogger = &CtxKey{Name: "logger"}

// WithLogger sets the logger to the context.
func WithLogger(ctx context.Context, logger logging.Logger) context.Context {
	return context.WithValue(ctx, CtxKeyLogger, logger)
}

// GetLogger gets the logger from the context.
// TODO: rename to FromContext and move to the loggingctx package.
//
//nolint:godox,ireturn // Necessary to fix the function in the interface.
func GetLogger(ctx context.Context) logging.Logger {
	value := ctx.Value(CtxKeyLogger)
	if value == nil {
		return nil
	}

	logger, ok := value.(logging.Logger)
	if !ok {
		return nil
	}

	return logger
}
