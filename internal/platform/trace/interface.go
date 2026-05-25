package trace

import (
	"context"
)

type ISpan interface {
	End()
	TraceID() string
	SpanID() string
	ParentSpanID() string
}

type ITracer interface {
	Start(ctx context.Context, spanName string) (context.Context, ISpan) // + opts
}

type ITracerProvider interface {
	Tracer(name string) ITracer // + opts
}

// Export to:
//   - Jaeger
//   - OTel Collector

// Один провайдер — много трейсеров: Функция initTracer вызывается строго один
// раз в main.go. А вот oteltrace.NewTracer("имя_пакета") вы можете создавать для
// каждого слоя приложения (для HTTP-хэндлеров, для репозиториев, для внешних клиентов),
// чтобы в панели мониторинга четко видеть, какой пакет создал спан.
