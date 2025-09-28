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
//   - logLevel: logging level;
//   - out: log output;
//   - format: log output format.
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
	config.TimeKey = "ts"
	config.LevelKey = "level"
	config.MessageKey = "msg"
	config.CallerKey = "caller"
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
		mapToZapCoreLevel(logLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zapSugarLogger := &ZapSugarLogger{
		log:      zapLogger.Sugar(),
		logLevel: logLevel,
	}

	zapSugarLogger.Info(
		"logger initialized",
		"level", zapSugarLogger.logLevel.String(),
		"format", format,
	)

	return zapSugarLogger, nil
}

// Debug logs the message and parameters with the debug level.
//
// Parameters:
//   - message: main log message;
//   - keysAndValues: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Debug(msg string, keysAndValues ...any) {
	if LevelDebug < l.logLevel {
		return
	}

	l.log.Debugw(msg, keysAndValues...)
}

// Info logs the message and parameters with the info level.
//
// Parameters:
//   - message: main log message;
//   - keysAndValues: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Info(msg string, keysAndValues ...any) {
	if LevelInfo < l.logLevel {
		return
	}

	l.log.Infow(msg, keysAndValues...)
}

// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
//
// Parameters:
//   - message: main log message;
//   - err: possible error;
//   - keysAndValues: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Warn(msg string, err error, keysAndValues ...any) {
	if LevelWarn < l.logLevel {
		return
	}

	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}

	l.log.Warnw(msg, keysAndValues...)
}

// Error logs a message and parameters with the error level and error.
//
// Parameters:
//   - message: main log message;
//   - err: error;
//   - keysAndValues: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Error(msg string, err error, keysAndValues ...any) {
	if LevelError < l.logLevel {
		return
	}

	keysAndValues = append(keysAndValues, "error", err.Error())

	l.log.Errorw(msg, keysAndValues...)
}

// Fatal logs a message and parameters with the fatal and critical error levels.
//
// Parameters:
//   - message: main log message;
//   - err: critical error;
//   - keysAndValues: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Fatal(msg string, err error, keysAndValues ...any) {
	if LevelFatal < l.logLevel {
		return
	}

	keysAndValues = append(keysAndValues, "error", err.Error())

	l.log.Fatalw(msg, keysAndValues...)
}

// Close releases resources used by the logger.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Close() error {
	_ = l.log.Sync() // На Windows zap.Sync может возвращать ошибку – игнорируем.

	return nil
}

// mapToZapCoreLevel — mapping LogLevel to zapcore.Level.
//
// Parameters:
//   - level: logging level.
func mapToZapCoreLevel(level LogLevel) zapcore.Level {
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
