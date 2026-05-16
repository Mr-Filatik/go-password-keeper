// Package tracectx provides utilities for embedding and retrieving trace information
// within a context.Context using a strongly-typed key.
package tracectx

import (
	"context"

	commonctx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx"
	"github.com/mr-filatik/go-password-keeper/internal/platform/trace"
)

//nolint:gochecknoglobals
var ctxTraceKey = &commonctx.CtxKey{Name: "trace"}

// SetTrace returns a new context containing the provided trace information.
func SetTrace(ctx context.Context, value trace.Trace) context.Context {
	return commonctx.SetValue(ctx, ctxTraceKey, value)
}

// GetTrace extracts the trace information from the context.
// It returns the trace data and a boolean indicating whether the trace was successfully found.
func GetTrace(ctx context.Context) (trace.Trace, bool) {
	return commonctx.GetValue[trace.Trace](ctx, ctxTraceKey)
}
