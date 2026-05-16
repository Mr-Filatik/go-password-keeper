package middleware

import (
	"context"
	"net/http"
	"strings"

	tracectx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/trace"
	"github.com/mr-filatik/go-password-keeper/internal/platform/trace"
)

const (
	HeaderTraceParent = "traceparent"

	HeaderTraceID = "X-Trace-Id"
	HeaderSpanID  = "X-Span-Id"
)

func AdvancedTracing() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := AdvancedTracingAction(r.Context(), r.Header)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdvancedTracingAction(ctx context.Context, header http.Header) context.Context {
	var traceInfo trace.Trace

	const (
		traceparentParts = 4
		// traceparentVersionIdx = 0
		traceparentTraceIdx  = 1
		traceparentParentIdx = 2
		// traceparentFlagsIdx = 3

		traceparentSeparator = "-"
	)

	// 1. Проверяем стандартный traceparent (формат: version-traceID-parentID-traceFlags)
	// Пример: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
	//
	//nolint:canonicalheader
	if tp := header.Get(HeaderTraceParent); tp != "" {
		parts := strings.Split(tp, traceparentSeparator)
		if len(parts) == traceparentParts {
			traceInfo = trace.New(parts[traceparentTraceIdx], parts[traceparentParentIdx])

			return tracectx.SetTrace(ctx, traceInfo)
		}
	}

	traceID := header.Get(HeaderTraceID)
	if traceID != "" {
		traceInfo = trace.New(traceID, header.Get(HeaderSpanID))

		return tracectx.SetTrace(ctx, traceInfo)
	}

	traceInfo = trace.NewEmpty()

	return tracectx.SetTrace(ctx, traceInfo)

	// w3cPropagator := propagation.TraceContext{}
	// ctx = w3cPropagator.Extract(ctx, propagation.HeaderCarrier(header))

	// // Если traceparent успешно распарсился и содержит валидный Span, возвращаем его
	// if trace.SpanContextFromContext(ctx).IsValid() {
	// 	return ctx
	// }

	// // 2. Если traceparent пустой/невалидный, ищем отдельные заголовки
	// traceIDStr := header.Get(HeaderTraceID)
	// spanIDStr := header.Get(HeaderSpanID)

	// if traceIDStr == "" {
	// 	return ctx // Ничего не нашли, возвращаем исходный контекст (OTel сам создаст новый span позже)
	// }

	// // Парсим строковые ID в бинарный формат OpenTelemetry
	// traceID, err := trace.TraceIDFromHex(traceIDStr)
	// if err != nil {
	// 	return ctx
	// }

	// // SpanID опционален, если его нет — создаем пустой/нулевой
	// var spanID trace.SpanID
	// if spanIDStr != "" {
	// 	if parsedSpanID, err := trace.SpanIDFromHex(spanIDStr); err == nil {
	// 		spanID = parsedSpanID
	// 	}
	// }

	// // Собираем валидный SpanContext
	// spanContext := trace.NewSpanContext(trace.SpanContextConfig{
	// 	TraceID:    traceID,
	// 	SpanID:     spanID,
	// 	TraceFlags: trace.FlagsSampled, // Помечаем, что трейс активен для сбора
	// 	Remote:     true,
	// })

	// // Записываем воссозданный Span в контекст
	// return trace.ContextWithSpanContext(ctx, spanContext)
}
