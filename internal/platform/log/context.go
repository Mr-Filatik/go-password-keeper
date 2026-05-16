// Package logging provides logging functionality.
package log

// Можно сделать ещё метод такой, чтобы caller был везде верный.
// Тогда надо запретить вызывать логгер напрямую. Но как, приватные поля?
func LogDebug(logger ILogger, msg string, options ...FieldOption) {
	if logger == nil {
		logger = Nop
	}

	logger.Debug(msg, options...)
}

// LogInfo writes a log to the logger located in the context with the info level.
func LogInfo(logger ILogger, msg string, options ...FieldOption) {
	if logger == nil {
		logger = Nop
	}

	logger.Info(msg, options...)
}

// LogWarn writes a log to the logger located in the context with the warning level.
func LogWarn(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = Nop
	}

	logger.Warn(msg, err, options...)
}

// LogError writes a log to the logger located in the context with the error level.
func LogError(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = Nop
	}

	logger.Error(msg, err, options...)
}

// LogFatal writes a log to the logger located in the context with the fatal level.
func LogFatal(logger ILogger, msg string, err error, options ...FieldOption) {
	if logger == nil {
		logger = Nop
	}

	logger.Fatal(msg, err, options...)
}
