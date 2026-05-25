package logwrapper

import (
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

type metrFunc func(level log.LogLevel)

type MetricWrapper struct {
	log.ILogger
	metrFn metrFunc
}

func NewMetricWrapper(logger log.ILogger, metrFn metrFunc) log.ILogger {
	if metrFn == nil {
		metrFn = func(_ log.LogLevel) {}
	}

	return &MetricWrapper{
		ILogger: logger,
		metrFn:  metrFn,
	}
}

func (l *MetricWrapper) With(options ...log.FieldOption) log.ILogger {
	return &MetricWrapper{
		ILogger: l.ILogger.With(options...),
		metrFn:  l.metrFn,
	}
}

func (l *MetricWrapper) Debug(msg string, options ...log.FieldOption) {
	l.ILogger.Debug(msg, options...)

	l.metrFn(log.LevelDebug)
}

func (l *MetricWrapper) Info(msg string, options ...log.FieldOption) {
	l.ILogger.Info(msg, options...)

	l.metrFn(log.LevelInfo)
}

func (l *MetricWrapper) Warn(msg string, err error, options ...log.FieldOption) {
	l.ILogger.Warn(msg, err, options...)

	l.metrFn(log.LevelWarn)
}

func (l *MetricWrapper) Error(msg string, err error, options ...log.FieldOption) {
	l.ILogger.Error(msg, err, options...)

	l.metrFn(log.LevelError)
}

func (l *MetricWrapper) Fatal(msg string, err error, options ...log.FieldOption) {
	l.ILogger.Fatal(msg, err, options...)

	l.metrFn(log.LevelFatal)
}

func (l *MetricWrapper) Close() error {
	return l.ILogger.Close()
}
