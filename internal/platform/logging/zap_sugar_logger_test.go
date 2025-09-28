package logging_test

import (
	"errors"
	"io"
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/mocks"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	TestNewZapSugarLogger
*/

type argsNewZapSugarLogger struct {
	logLevel logging.LogLevel
	format   logging.LogFormat
}

type wantNewZapSugarLogger struct {
	data []string
}

type testNewZapSugarLogger struct {
	name string
	args argsNewZapSugarLogger
	want wantNewZapSugarLogger
}

func getTestsNewZapSugarLogger() []testNewZapSugarLogger {
	getDataInJSON := func(loggerLogLevel logging.LogLevel) []string {
		return []string{
			"\"level\":\"INFO\"",
			"\"ts\":",
			"\"caller\":\"logging/zap_sugar_logger.go:", // without line
			"\"msg\":\"logger initialized\"",
			"\"level\":\"" + loggerLogLevel.String() + "\"",
			"\"format\":\"JSON\"",
		}
	}

	return []testNewZapSugarLogger{
		{
			name: "debug level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelDebug,
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: getDataInJSON(logging.LevelDebug),
			},
		},
		{
			name: "info level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelInfo,
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: getDataInJSON(logging.LevelInfo),
			},
		},
		{
			name: "info level in text format",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelInfo,
				format:   logging.FormatText,
			},
			want: wantNewZapSugarLogger{
				data: []string{
					"INFO",
					"logging/zap_sugar_logger.go:", // without line
					"logger initialized",
					"\"level\": \"info\"",
					"\"format\": \"TEXT\"",
				},
			},
		},
		{
			name: "info level in unknown format",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelInfo,
				format:   logging.LogFormat("UNKNOWN"),
			},
			want: wantNewZapSugarLogger{
				data: getDataInJSON(logging.LevelInfo),
			},
		},
		{
			name: "warning level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelWarn,
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: nil,
			},
		},
		{
			name: "error level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelError,
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: nil,
			},
		},
		{
			name: "fatal level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LevelFatal,
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: nil,
			},
		},
		{
			name: "unknown level",
			args: argsNewZapSugarLogger{
				logLevel: logging.LogLevel(99),
				format:   logging.FormatJSON,
			},
			want: wantNewZapSugarLogger{
				data: nil,
			},
		},
	}
}

func TestNewZapSugarLogger(t *testing.T) {
	t.Parallel()

	tests := getTestsNewZapSugarLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := mocks.NewMockWriter()

			require.NotNilf(t, mockWriter, "NewMockWriter() writer = nil")

			logger, err := logging.NewZapSugarLogger(tt.args.logLevel, mockWriter, tt.args.format)

			require.NoErrorf(t, err, "NewZapSugarLogger() error = %v, want nil", err)
			require.NotNilf(t, logger, "NewZapSugarLogger() logger = nil")

			if tt.want.data != nil {
				lastLog, _ := mockWriter.GetUnreadedData()
				for _, item := range tt.want.data {
					assert.Containsf(t, string(lastLog), item, "last log not contains %v", item)
				}
			}
		})
	}
}

/*
	TestZapSugarLogger
*/

var errTestReason = errors.New("test reason")

type argsZapSugarLogger struct {
	msg           string
	err           error
	keysAndValues []any
}

type wantZapSugarLogger struct {
	data []string
}

type testZapSugarLogger struct {
	name           string
	loggerLogLevel logging.LogLevel
	args           argsZapSugarLogger
	want           wantZapSugarLogger
}

/*
	TestZapSugarLogger_Debug
*/

func getTestsZapSugarLoggerDebug() []testZapSugarLogger {
	return []testZapSugarLogger{
		{
			name:           "debug logger",
			loggerLogLevel: logging.LevelDebug,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"DEBUG\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
				},
			},
		},
		{
			name:           "info logger",
			loggerLogLevel: logging.LevelInfo,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "warning logger",
			loggerLogLevel: logging.LevelWarn,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "error logger",
			loggerLogLevel: logging.LevelError,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "fatal logger",
			loggerLogLevel: logging.LevelFatal,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
	}
}

func TestZapSugarLogger_Debug(t *testing.T) {
	t.Parallel()

	logLevel := logging.LevelDebug
	tests := getTestsZapSugarLoggerDebug()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := createMockWriter(t)
			logger := createLogger(t, mockWriter, tt.loggerLogLevel)

			mockWriter.MarkDataAsRead()

			logger.Debug(tt.args.msg, tt.args.keysAndValues...)

			if tt.want.data != nil {
				t.Logf("\n[%s]\nthe %s level log was written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)

				lastLog, isLastLog := mockWriter.GetUnreadedData()

				require.Truef(t, isLastLog, "the log should have been written, but it is missing")

				for _, item := range tt.want.data {
					assert.Containsf(t, string(lastLog), item, "last log not contains %v", item)
				}
			} else {
				t.Logf("\n[%s]\nthe %s level log was NOT written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)
			}
		})
	}
}

/*
	TestZapSugarLogger_Info
*/

func getTestsZapSugarLoggerInfo() []testZapSugarLogger {
	return []testZapSugarLogger{
		{
			name:           "debug logger",
			loggerLogLevel: logging.LevelDebug,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"INFO\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
				},
			},
		},
		{
			name:           "info logger",
			loggerLogLevel: logging.LevelInfo,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"INFO\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
				},
			},
		},
		{
			name:           "warning logger",
			loggerLogLevel: logging.LevelWarn,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "error logger",
			loggerLogLevel: logging.LevelError,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "fatal logger",
			loggerLogLevel: logging.LevelFatal,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
	}
}

func TestZapSugarLogger_Info(t *testing.T) {
	t.Parallel()

	logLevel := logging.LevelInfo
	tests := getTestsZapSugarLoggerInfo()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := createMockWriter(t)
			logger := createLogger(t, mockWriter, tt.loggerLogLevel)

			mockWriter.MarkDataAsRead()

			logger.Info(tt.args.msg, tt.args.keysAndValues...)

			if tt.want.data != nil {
				t.Logf("\n[%s]\nthe %s level log was written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)

				lastLog, isLastLog := mockWriter.GetUnreadedData()

				require.Truef(t, isLastLog, "the log should have been written, but it is missing")

				for _, item := range tt.want.data {
					assert.Containsf(t, string(lastLog), item, "last log not contains %v", item)
				}
			} else {
				t.Logf("\n[%s]\nthe %s level log was NOT written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)
			}
		})
	}
}

/*
	TestZapSugarLogger_Warn
*/

func getTestsZapSugarLoggerWarn() []testZapSugarLogger {
	return []testZapSugarLogger{
		{
			name:           "debug logger",
			loggerLogLevel: logging.LevelDebug,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"WARN\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "info logger",
			loggerLogLevel: logging.LevelInfo,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"WARN\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "warning logger",
			loggerLogLevel: logging.LevelWarn,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"WARN\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "warning logger with error",
			loggerLogLevel: logging.LevelWarn,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           nil,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"WARN\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
				},
			},
		},
		{
			name:           "error logger",
			loggerLogLevel: logging.LevelError,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
		{
			name:           "fatal logger",
			loggerLogLevel: logging.LevelFatal,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
	}
}

func TestZapSugarLogger_Warn(t *testing.T) {
	t.Parallel()

	logLevel := logging.LevelWarn
	tests := getTestsZapSugarLoggerWarn()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := createMockWriter(t)
			logger := createLogger(t, mockWriter, tt.loggerLogLevel)

			mockWriter.MarkDataAsRead()

			logger.Warn(tt.args.msg, tt.args.err, tt.args.keysAndValues...)

			if tt.want.data != nil {
				t.Logf("\n[%s]\nthe %s level log was written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)

				lastLog, isLastLog := mockWriter.GetUnreadedData()

				require.Truef(t, isLastLog, "the log should have been written, but it is missing")

				for _, item := range tt.want.data {
					assert.Containsf(t, string(lastLog), item, "last log not contains %v", item)
				}
			} else {
				t.Logf("\n[%s]\nthe %s level log was NOT written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)
			}
		})
	}
}

/*
	TestZapSugarLogger_Error
*/

func getTestsZapSugarLoggerError() []testZapSugarLogger {
	return []testZapSugarLogger{
		{
			name:           "debug logger",
			loggerLogLevel: logging.LevelDebug,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"ERROR\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "info logger",
			loggerLogLevel: logging.LevelInfo,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"ERROR\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "warning logger",
			loggerLogLevel: logging.LevelWarn,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"ERROR\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "error logger",
			loggerLogLevel: logging.LevelError,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: []string{
					"\"level\":\"ERROR\"",
					"\"caller\":\"logging/zap_sugar_logger_test.go:", // without line
					"\"msg\":\"test message\"",
					"\"test key\":\"test value\"",
					"\"error\":\"test reason\"",
				},
			},
		},
		{
			name:           "fatal logger",
			loggerLogLevel: logging.LevelFatal,
			args: argsZapSugarLogger{
				msg:           "test message",
				err:           errTestReason,
				keysAndValues: []any{"test key", "test value"},
			},
			want: wantZapSugarLogger{
				data: nil,
			},
		},
	}
}

func TestZapSugarLogger_Error(t *testing.T) {
	t.Parallel()

	logLevel := logging.LevelError
	tests := getTestsZapSugarLoggerError()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := createMockWriter(t)
			logger := createLogger(t, mockWriter, tt.loggerLogLevel)

			mockWriter.MarkDataAsRead()

			logger.Error(tt.args.msg, tt.args.err, tt.args.keysAndValues...)

			if tt.want.data != nil {
				t.Logf("\n[%s]\nthe %s level log was written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)

				lastLog, isLastLog := mockWriter.GetUnreadedData()

				require.Truef(t, isLastLog, "the log should have been written, but it is missing")

				for _, item := range tt.want.data {
					assert.Containsf(t, string(lastLog), item, "last log not contains %v", item)
				}
			} else {
				t.Logf("\n[%s]\nthe %s level log was NOT written by the %s level logger",
					tt.name, logLevel, tt.loggerLogLevel)
			}
		})
	}
}

/*
	TestZapSugarLogger_Fatal
*/

func TestZapSugarLogger_Fatal(t *testing.T) {
	t.Parallel()

	// zap.SugaredLogger calls os.Exit(1) at Fatal level
	// For tests, call WithOptions(zap.WithFatalHook(zapcore.WriteThenNoop))
}

/*
	TestZapSugarLogger_Close
*/

func TestZapSugarLogger_Close(t *testing.T) {
	t.Parallel()

	mockWriter := createMockWriter(t)

	logger := createLogger(t, mockWriter, logging.LevelInfo)

	err := logger.Close()

	require.NoErrorf(t, err, "error closing logger")
}

/*
	Helpers
*/

func createMockWriter(t *testing.T) *mocks.MockWriter {
	t.Helper()

	mockWriter := mocks.NewMockWriter()

	require.NotNilf(t, mockWriter, "NewMockWriter() return nil")

	return mockWriter
}

func createLogger(
	t *testing.T,
	mockWriter io.Writer,
	loggerLogLevel logging.LogLevel,
) *logging.ZapSugarLogger {
	t.Helper()

	logger, err := logging.NewZapSugarLogger(loggerLogLevel, mockWriter, logging.FormatJSON)

	require.NoErrorf(t, err, "NewZapSugarLogger() error = %v, want nil", err)
	require.NotNilf(t, logger, "NewZapSugarLogger() logger = nil")

	return logger
}
