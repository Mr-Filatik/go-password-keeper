package logctx

import (
	"context"

	tracectx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/trace"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// Debug logs a message at debug level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Debug(ctx context.Context, msg string, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Debug(msg, options...)
}

// Info logs a message at info level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Info(ctx context.Context, msg string, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Info(msg, options...)
}

// Warn logs a message and an error at warn level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Warn(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Warn(msg, err, options...)
}

// Error logs a message and an error at error level.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Error(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Error(msg, err, options...)
}

// Fatal logs a message and an error at fatal level, then terminates the process.
// It automatically appends trace information from the context if auto-tracing is enabled.
func Fatal(ctx context.Context, msg string, err error, options ...log.FieldOption) {
	logger, options := prepareLoggerAndOptions(ctx, options...)
	logger.Fatal(msg, err, options...)
}

//nolint:ireturn
func prepareLoggerAndOptions(ctx context.Context, options ...log.FieldOption) (log.ILogger, []log.FieldOption) {
	logger := GetLogger(ctx)

	// autoLogger, ok := logger.(log.IAutoLogger)
	// if ok && autoLogger.IsAutoTracing() {
	// 	// omit

	// 	trace, ok := tracectx.GetTrace(ctx)
	// 	if ok {
	// 		options = append(options, log.WithAdvancedTraceField(trace)...)
	// 	}
	// }

	trace, ok := tracectx.GetTrace(ctx)
	if !ok {
		return logger, options
	}

	options = append(options,
		log.WithTraceIDField(trace.TraceID()),
		log.WithSpanIDField(trace.SpanID()),
	)

	if trace.ParentSpanID() != "" {
		options = append(options, log.WithParentIDField(trace.ParentSpanID()))
	}

	// span := trace.SpanFromContext(ctx)
	// spanContext := span.SpanContext()

	// if !spanContext.IsValid() {
	// 	return logger, options
	// }

	// options = append(options,
	// 	log.WithTraceIDField(spanContext.TraceID().String()),
	// 	log.WithSpanIDField(spanContext.SpanID().String()),
	// )

	// if readOnlySpan, ok := span.(sdktrace.ReadOnlySpan); ok {
	// 	// Получаем структуру родителя
	// 	parent := readOnlySpan.Parent()

	// 	// Проверяем, что у спана вообще был родитель
	// 	if parent.IsValid() {
	// 		options = append(options, log.WithParentIDField(parent.SpanID().String()))
	// 	}
	// }

	return logger, options
}

// https://github.com/uptrace/opentelemetry-go-extra
