// Package logging provides logging functionality.
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapSugarLogger - adapter for *zap.SugaredLogger.
type ZapSugarLogger struct {
	// log - Logger.
	log *zap.SugaredLogger

	// logLevel - Logging level.
	logLevel LogLevel

	globalFields []any
}

// NewZapSugarLogger creates a new *ZapSugarLogger logger instance.
//
// Parameters:
//   - logLevel LogLevel: logging level;
//   - out io.Writer: log output;
//   - format LogFormat: log output format.
func NewZapSugarLogger(
	logLevel LogLevel,
	options ...ConfigOption,
) (Logger, error) {
	config := defaultConfig()

	for _, option := range options {
		option(&config)
	}

	logLevel = logLevel.Validate()

	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.TimeKey = FieldBaseTimestamp
	zapConfig.LevelKey = FieldBaseLevel
	zapConfig.MessageKey = FieldBaseMessage
	zapConfig.CallerKey = FieldBaseCaller
	zapConfig.EncodeLevel = zapcore.CapitalLevelEncoder // for text can be replaced with CapitalColorLevelEncoder

	var encoder zapcore.Encoder

	switch config.format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(zapConfig)
	case FormatText:
		encoder = zapcore.NewConsoleEncoder(zapConfig)
	default:
		encoder = zapcore.NewJSONEncoder(zapConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(config.writer),
		levelToZapCoreLevel(logLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(config.callerSkipCount))
	zapSugarLogger := &ZapSugarLogger{
		log:          zapLogger.Sugar(),
		logLevel:     logLevel,
		globalFields: config.globalFields,
	}

	LogInfo(zapSugarLogger, "ZapSugar logger initialize is successful",
		WithDataField(map[string]any{
			"format": config.format,
			"level":  logLevel.String(),
		}))

	return zapSugarLogger, nil
}

// With returns a new logger with added common fields (labels).
//
// Parameters:
//   - keysAndValues ...any: common fields (labels).
//
// Implements the internal/platform/logging.Logger interface.
//
//nolint:ireturn // Necessary to fix the function in the interface.
func (l *ZapSugarLogger) With(options ...FieldOption) Logger {
	if l == nil {
		return nil
	}

	// TODO убрать дубли
	fields := l.globalFields
	fields = append(fields, applyOptions(options...)...)

	return &ZapSugarLogger{
		log:          l.log.With(fields...),
		logLevel:     l.logLevel,
		globalFields: fields,
	}
}

// Debug logs the message and parameters with the debug level.
//
// Parameters:
//   - msg string: main log message;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Debug(msg string, options ...FieldOption) {
	if LevelDebug < l.logLevel {
		return
	}

	fields := l.globalFields
	fields = append(fields, applyOptions(options...)...)

	l.log.Debugw(msg, fields...)
}

// Info logs the message and parameters with the info level.
//
// Parameters:
//   - msg string: main log message;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Info(msg string, options ...FieldOption) {
	if LevelInfo < l.logLevel {
		return
	}

	fields := l.globalFields
	fields = append(fields, applyOptions(options...)...)

	l.log.Infow(msg, fields...)
}

// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
//
// Parameters:
//   - msg string: main log message;
//   - err error: possible error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Warn(msg string, err error, options ...FieldOption) {
	if LevelWarn < l.logLevel {
		return
	}

	var allOptions []FieldOption
	if err != nil {
		allOptions = append([]FieldOption{WithErrorField(err)}, options...)
	} else {
		allOptions = options
	}

	fields := l.globalFields
	fields = append(fields, applyOptions(allOptions...)...)

	l.log.Warnw(msg, fields...)
}

// Error logs a message and parameters with the error level and error.
//
// Parameters:
//   - msg string: main log message;
//   - err error: error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Error(msg string, err error, options ...FieldOption) {
	if LevelError < l.logLevel {
		return
	}

	allOptions := append([]FieldOption{WithErrorField(err)}, options...)

	fields := l.globalFields
	fields = append(fields, applyOptions(allOptions...)...)

	l.log.Errorw(msg, fields...)
}

// Fatal logs a message and parameters with the fatal and critical error levels.
//
// Parameters:
//   - msg string: main log message;
//   - err error: critical error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Fatal(msg string, err error, options ...FieldOption) {
	if LevelFatal < l.logLevel {
		return
	}

	allOptions := append([]FieldOption{WithErrorField(err)}, options...)

	fields := l.globalFields
	fields = append(fields, applyOptions(allOptions...)...)

	l.log.Fatalw(msg, fields...)
}

// Close releases resources used by the logger.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Close() error {
	_ = l.log.Sync() // На Windows zap.Sync может возвращать ошибку – игнорируем.

	return nil
}

// levelToZapCoreLevel — mapping LogLevel to zapcore.Level.
//
// Parameters:
//   - level LogLevel: logging level.
func levelToZapCoreLevel(level LogLevel) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.ErrorLevel
	}
}
