package oteltrace

import (
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// OTelErrorHandler реализует интерфейс otel.ErrorHandler
type OTelErrorHandler struct {
	logger log.ILogger
}

func NewOTelErrorHandler(myLogger log.ILogger) *OTelErrorHandler {
	return &OTelErrorHandler{
		// Добавляем маркер, чтобы в JSON было видно, что это сбой инфраструктуры
		logger: myLogger.With(log.WithComponentField("opentelemetry-package")),
	}
}

func (h *OTelErrorHandler) Handle(err error) {
	// Вызываем ваш метод Error, передавая ошибку инфраструктуры
	h.logger.Error("OpenTelemetry критический сбой подсистемы", err)
}
