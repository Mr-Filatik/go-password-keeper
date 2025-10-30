package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func DefaultConveyour() []Middleware {
	return []Middleware{}
}

// [Recover]           // самый внешний; формирует 500, тело/заголовки при панике
// [RequestID]         // чтобы ID попал в логи/метрики
// [Logging]           // при панике логирует и re-panic
// [Metrics]           // при панике пишет метрику 500 и re-panic
// [Auth/RateLimit/...]
// [Handlers]
