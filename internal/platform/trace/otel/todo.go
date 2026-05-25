package oteltrace

// package oteltrace

// import (
// 	"context"

// 	mytrace "github.com/mr-filatik/go-password-keeper/internal/platform/trace"
// 	"go.opentelemetry.io/otel/trace"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// )

// type Tracer struct {
// 	tr trace.Tracer
// }

// func (t *Tracer) Start(ctx context.Context, name string) (context.Context, mytrace.ISpan) {
// 	// 1. До старта нового спана смотрим на контекст.
// 	// Если там уже есть валидный спан — он является РОДИТЕЛЕМ для создаваемого.
// 	parentCtx := trace.SpanFromContext(ctx).SpanContext()
// 	var parentID string
// 	if parentCtx.IsValid() {
// 		parentID = parentCtx.SpanID().String()
// 	}

// 	// 2. Стартуем новый спан стандартным образом
// 	newCtx, otelSpan := t.tr.Start(ctx, name)

// 	// 3. Возвращаем новый контекст и наш адаптер спана,
// 	// принудительно передавая туда вычисленный parentID
// 	mySpan := Span{
// 		Span:         otelSpan,
// 		parentSpanID: parentID,
// 	}

// 	return newCtx, mySpan
// }

// // ==========================================

// type Span struct {
// 	trace.Span
// 	parentSpanID string // Храним вычисленный ID родителя прямо в структуре
// }

// func (s Span) End() {
// 	s.Span.End()
// }

// func (s Span) TraceID() string {
// 	return s.Span.SpanContext().TraceID().String()
// }

// func (s Span) SpanID() string {
// 	return s.Span.SpanContext().SpanID().String()
// }

// func (s Span) ParentSpanID() string {
// 	// Если мы зафиксировали родителя при старте, возвращаем его
// 	if s.parentSpanID != "" {
// 		return s.parentSpanID
// 	}

// 	// Запасной вариант для локальных спанов (на случай, если спан создавался в обход метода Start)
// 	if readOnlySpan, ok := s.Span.(sdktrace.ReadOnlySpan); ok {
// 		if parent := readOnlySpan.Parent(); parent.IsValid() {
// 			return parent.SpanID().String()
// 		}
// 	}

// 	return ""
// }
