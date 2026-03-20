// Package logging provides logging functionality.
package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapSugarLogger - adapter for *zap.SugaredLogger.
type ZapSugarLogger struct {
	// log - Logger.
	log *zap.SugaredLogger

	// logLevel - Logging level.
	logLevel LogLevel
}

// NewZapSugarLogger creates a new *ZapSugarLogger logger instance.
//
// Parameters:
//   - logLevel LogLevel: logging level;
//   - out io.Writer: log output;
//   - format LogFormat: log output format.
func NewZapSugarLogger(
	logLevel LogLevel,
	out io.Writer,
	format LogFormat,
) (*ZapSugarLogger, error) {
	if out == nil {
		out = os.Stdout
	}

	logLevel = logLevel.Validate()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = FieldBaseTimestamp
	config.LevelKey = FieldBaseLevel
	config.MessageKey = FieldBaseMessage
	config.CallerKey = FieldBaseCaller
	config.EncodeLevel = zapcore.CapitalLevelEncoder // for text can be replaced with CapitalColorLevelEncoder

	var encoder zapcore.Encoder

	switch format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(config)
	case FormatText:
		encoder = zapcore.NewConsoleEncoder(config)
	default:
		format = FormatJSON
		encoder = zapcore.NewJSONEncoder(config)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(out),
		levelToZapCoreLevel(logLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zapSugarLogger := &ZapSugarLogger{
		log:      zapLogger.Sugar(),
		logLevel: logLevel,
	}

	zapSugarLogger.Info("ZapSugar logger initialize is successful",
		WithCustomField("format", format),
		WithCustomField("level", logLevel.String()),
	)

	return zapSugarLogger, nil
}

// NewZapSugarLoggerWithFields creates a new *ZapSugarLogger logger instance with common fields (labels).
//
// Parameters:
//   - logLevel LogLevel: logging level;
//   - out io.Writer: log output;
//   - format LogFormat: log output format;
//   - keysAndValues ...any: common fields (labels).
//
//nolint:ireturn // Necessary to fix the function in the interface.
func NewZapSugarLoggerWithFields(
	logLevel LogLevel,
	out io.Writer,
	format LogFormat,
	keysAndValues ...any,
) (Logger, error) {
	logger, err := NewZapSugarLogger(logLevel, out, format)
	if err != nil {
		return logger, err
	}

	return logger.With(keysAndValues...), nil
}

// With returns a new logger with added common fields (labels).
//
// Parameters:
//   - keysAndValues ...any: common fields (labels).
//
// Implements the internal/platform/logging.Logger interface.
//
//nolint:ireturn // Necessary to fix the function in the interface.
func (l *ZapSugarLogger) With(keysAndValues ...any) Logger {
	if l == nil {
		return nil
	}

	return &ZapSugarLogger{
		log:      l.log.With(keysAndValues...),
		logLevel: l.logLevel,
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

	fields := applyOptions(options...)

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

	fields := applyOptions(options...)

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

	fields := applyOptions(allOptions...)

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
	fields := applyOptions(allOptions...)

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
	fields := applyOptions(allOptions...)

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
