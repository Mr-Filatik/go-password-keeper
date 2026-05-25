package oteltrace

import (
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	trace.Span

	// traceID      string
	// spanID       string
	parentSpanID string
}

func (s Span) End() {
	s.Span.End()
}

func (s Span) TraceID() string {
	// return s.traceID

	return s.Span.SpanContext().TraceID().String() // or [16]byte
}

func (s Span) SpanID() string {
	// return s.spanID

	return s.Span.SpanContext().SpanID().String()
}

func (s Span) ParentSpanID() string {
	return s.parentSpanID

	// if readOnlySpan, ok := s.Span.(sdktrace.ReadOnlySpan); ok {
	// 	parentContext := readOnlySpan.Parent()
	// 	if parentContext.IsValid() {
	// 		return parentContext.SpanID().String()
	// 	}
	// }

	// return ""
}
