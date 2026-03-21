// Package zaplog provides logging functionality.
package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// ZapSugarLogger - adapter for *zap.SugaredLogger.
type ZapSugarLogger struct {
	// log - Logger.
	log *zap.SugaredLogger

	// logLevel - Logging level.
	logLevel log.LogLevel

	globalFields []any
}

// NewZapSugarLogger creates a new *ZapSugarLogger logger instance.
//
// Parameters:
//   - logLevel LogLevel: logging level;
//   - out io.Writer: log output;
//   - format LogFormat: log output format.
func NewZapSugarLogger(
	logLevel log.LogLevel,
	options ...log.ConfigOption,
) (log.ILogger, error) {
	config := log.DefaultConfig()

	for _, option := range options {
		option(&config)
	}

	logLevel = logLevel.Validate()

	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.TimeKey = log.FieldBaseTimestamp
	zapConfig.LevelKey = log.FieldBaseLevel
	zapConfig.MessageKey = log.FieldBaseMessage
	zapConfig.CallerKey = log.FieldBaseCaller
	zapConfig.EncodeLevel = zapcore.CapitalLevelEncoder // for text can be replaced with CapitalColorLevelEncoder

	var encoder zapcore.Encoder

	switch config.GetFormat() {
	case log.FormatJSON:
		encoder = zapcore.NewJSONEncoder(zapConfig)
	case log.FormatText:
		encoder = zapcore.NewConsoleEncoder(zapConfig)
	default:
		encoder = zapcore.NewJSONEncoder(zapConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(config.GetWriter()),
		levelToZapCoreLevel(logLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(config.GetCallerSkipCount())).Sugar()
	zapSugarLogger := &ZapSugarLogger{
		log:          zapLogger.With(config.GetGlobalFields()...),
		logLevel:     logLevel,
		globalFields: config.GetGlobalFields(),
	}

	log.LogInfo(zapSugarLogger, "ZapSugar logger initialize is successful",
		log.WithDataField(map[string]any{
			"format": config.GetFormat(),
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
func (l *ZapSugarLogger) With(options ...log.FieldOption) log.ILogger {
	if l == nil {
		return nil
	}

	newOpts := log.ApplyOptions(options...)
	if len(newOpts) == 0 {
		return l
	}

	totalCap := len(l.globalFields) + len(newOpts)
	unique := make([]any, 0, totalCap)

	seen := make(map[any]struct{}, totalCap/2)

	for i := 0; i < len(l.globalFields); i += 2 {
		key := l.globalFields[i]
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
		}
	}

	for i := 0; i < len(newOpts); i += 2 {
		key := newOpts[i]
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			unique = append(unique, key, newOpts[i+1])
		}
	}

	return &ZapSugarLogger{
		log:          l.log.With(unique...),
		logLevel:     l.logLevel,
		globalFields: unique,
	}
}

// Debug logs the message and parameters with the debug level.
//
// Parameters:
//   - msg string: main log message;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Debug(msg string, options ...log.FieldOption) {
	if log.LevelDebug < l.logLevel {
		return
	}

	l.log.Debugw(msg, log.ApplyOptions(options...)...)
}

// Info logs the message and parameters with the info level.
//
// Parameters:
//   - msg string: main log message;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Info(msg string, options ...log.FieldOption) {
	if log.LevelInfo < l.logLevel {
		return
	}

	l.log.Infow(msg, log.ApplyOptions(options...)...)
}

// Warn logs a message and parameters with the warn level and a possible (non-critical) error.
//
// Parameters:
//   - msg string: main log message;
//   - err error: possible error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Warn(msg string, err error, options ...log.FieldOption) {
	if log.LevelWarn < l.logLevel {
		return
	}

	var allOptions []log.FieldOption
	if err != nil {
		allOptions = append([]log.FieldOption{log.WithErrorField(err)}, options...)
	} else {
		allOptions = options
	}

	l.log.Warnw(msg, log.ApplyOptions(allOptions...)...)
}

// Error logs a message and parameters with the error level and error.
//
// Parameters:
//   - msg string: main log message;
//   - err error: error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Error(msg string, err error, options ...log.FieldOption) {
	if log.LevelError < l.logLevel {
		return
	}

	allOptions := append([]log.FieldOption{log.WithErrorField(err)}, options...)

	l.log.Errorw(msg, log.ApplyOptions(allOptions...)...)
}

// Fatal logs a message and parameters with the fatal and critical error levels.
//
// Parameters:
//   - msg string: main log message;
//   - err error: critical error;
//   - datas ...any: additional information as a key-value pair.
//
// Implements the internal/platform/logging.Logger interface.
func (l *ZapSugarLogger) Fatal(msg string, err error, options ...log.FieldOption) {
	if log.LevelFatal < l.logLevel {
		return
	}

	allOptions := append([]log.FieldOption{log.WithErrorField(err)}, options...)

	l.log.Fatalw(msg, log.ApplyOptions(allOptions...)...)
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
func levelToZapCoreLevel(level log.LogLevel) zapcore.Level {
	switch level {
	case log.LevelDebug:
		return zapcore.DebugLevel
	case log.LevelInfo:
		return zapcore.InfoLevel
	case log.LevelWarn:
		return zapcore.WarnLevel
	case log.LevelError:
		return zapcore.ErrorLevel
	case log.LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.ErrorLevel
	}
}
