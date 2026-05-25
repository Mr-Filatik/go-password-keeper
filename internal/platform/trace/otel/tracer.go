package oteltrace

import (
	"context"

	tracectx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/trace"
	mytrace "github.com/mr-filatik/go-password-keeper/internal/platform/trace"
	"go.opentelemetry.io/otel/trace"
)

var otr mytrace.ITracer = &Tracer{}

type Tracer struct {
	trace.Tracer
}

func (t *Tracer) Start(ctx context.Context, spanName string) (context.Context, mytrace.ISpan) {
	parentCtx := trace.SpanFromContext(ctx).SpanContext()
	var parentID string
	if parentCtx.IsValid() {
		parentID = parentCtx.SpanID().String()
	}

	newCtx, otelSpan := t.Tracer.Start(ctx, spanName)
	// newSpanCtx := otelSpan.SpanContext()

	span := &Span{
		Span: otelSpan,
		// TraceID: newSpanCtx.TraceID().String(),
		// SpanID: newSpanCtx.SpanID().String(),
		parentSpanID: parentID,
	}

	newCtx = tracectx.SetTrace(newCtx, span)
	//tracectx.SetTrace(newCtx, span), span

	return newCtx, span
}
