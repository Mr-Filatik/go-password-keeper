package middleware

import (
	"net/http"

	"github.com/mr-filatik/go-password-keeper/internal/platform/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func Tracing(tracer trace.ITracer) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := tracer.Start(ctx, "tracing-middleware")
			defer span.End()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
