package oteltrace

import (
	"github.com/go-logr/logr"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// OTelLoggerAdapter переводит вызовы из OTel (logr.LogSink) в ваш ILogger
type OTelLoggerAdapter struct {
	logger log.ILogger
}

// NewOTelLoggerAdapter оборачивает ваш кастомный ILogger для передачи в OpenTelemetry
func NewOTelLoggerAdapter(myLogger log.ILogger) logr.Logger {
	sink := &OTelLoggerAdapter{
		// Сразу добавляем метку компонента через ваш метод With
		// (Вместо trace.Field используйте вашу фабрику FieldOption, например logger.String("...", "..."))
		logger: myLogger.With(log.WithComponentField("opentelemetry-package")),
	}
	return logr.New(sink)
}

func (o *OTelLoggerAdapter) Init(info logr.RuntimeInfo) {}

func (o *OTelLoggerAdapter) Enabled(level int) bool {
	return true
}

// Info перехватывает информационные логи OTel
func (o *OTelLoggerAdapter) Info(level int, msg string, keysAndValues ...interface{}) {
	options := o.toFieldOptions(keysAndValues)

	// Маппинг уровней OTel на методы вашего ILogger
	switch {
	case level >= 8:
		o.logger.Debug(msg, options...)
	case level >= 4:
		o.logger.Info(msg, options...)
	default:
		// В вашем интерфейсе Warn требует ошибку первым аргументом,
		// передаем nil, так как это просто информационное предупреждение
		o.logger.Warn(msg, nil, options...)
	}
}

// Error перехватывает внутренние сбои OTel (например, проблемы с Jaeger)
func (o *OTelLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	options := o.toFieldOptions(keysAndValues)

	// Вызываем ваш метод Error, который идеально принимает сообщение и ошибку
	o.logger.Error(msg, err, options...)
}

// toFieldOptions конвертирует сырой слайс OTel (ключ-значение) в ваши FieldOption
func (o *OTelLoggerAdapter) toFieldOptions(kv []interface{}) []log.FieldOption {
	if len(kv) == 0 {
		return nil
	}

	options := make([]log.FieldOption, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			if key, ok := kv[i].(string); ok {
				// ВАЖНО: Используйте здесь вашу реальную функцию создания FieldOption.
				// Например: logger.Any(key, kv[i+1])
				options = append(options, log.WithCustomField(key, kv[i+1]))
			}
		}
	}
	return options
}

// WithValues добавляет поля в логгер (вызывается внутренностями OTel)
func (o *OTelLoggerAdapter) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &OTelLoggerAdapter{
		logger: o.logger.With(o.toFieldOptions(keysAndValues)...),
	}
}

// WithName добавляет имя логгера OTel (например, "otlp-trace-exporter")
func (o *OTelLoggerAdapter) WithName(name string) logr.LogSink {
	return &OTelLoggerAdapter{
		logger: o.logger.With(log.WithCustomField("otel_logger_name", name)),
	}
}
