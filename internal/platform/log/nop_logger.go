// Package log provides logging functionality.
package log

type nopLogger struct{}

var Nop ILogger = &nopLogger{}

func (l *nopLogger) With(_ ...FieldOption) ILogger {
	return l
}

func (l *nopLogger) Debug(_ string, _ ...FieldOption) {}

func (l *nopLogger) Info(_ string, _ ...FieldOption) {}

func (l *nopLogger) Warn(_ string, _ error, _ ...FieldOption) {}

func (l *nopLogger) Error(_ string, _ error, _ ...FieldOption) {}

func (l *nopLogger) Fatal(_ string, _ error, _ ...FieldOption) {}

func (l *nopLogger) Close() error {
	return nil
}
