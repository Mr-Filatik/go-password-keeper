package tracectx

import (
	"context"

	"github.com/mr-filatik/go-password-keeper/internal/platform/trace"
)

// Next advances the trace context to the next step, ignoring the returned trace object.
// If no trace exists in the context, a new empty trace is initialized.
func Next(ctx context.Context) context.Context {
	ctx, _ = NextWithTrace(ctx)

	return ctx
}

// NextWithTrace advances the trace context to the next step and returns both the updated context
// and the new trace object.
// If no trace is present in the context, it generates a new empty trace, attaches it, and returns it.
func NextWithTrace(ctx context.Context) (context.Context, trace.Trace) {
	currentTrace, ok := GetTrace(ctx)
	if !ok {
		newTrace := trace.NewEmpty()

		return SetTrace(ctx, newTrace), newTrace
	}

	newTrace := currentTrace.Next()

	return SetTrace(ctx, newTrace), newTrace
}
