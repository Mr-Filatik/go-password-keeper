package logging_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    logging.LogLevel
		want string
	}{
		{
			name: "debug level",
			l:    logging.LevelDebug,
			want: "debug",
		},
		{
			name: "info level",
			l:    logging.LevelInfo,
			want: "info",
		},
		{
			name: "warning level",
			l:    logging.LevelWarn,
			want: "warn",
		},
		{
			name: "error level",
			l:    logging.LevelError,
			want: "error",
		},
		{
			name: "fatal level",
			l:    logging.LevelFatal,
			want: "fatal",
		},
		{
			name: "unknown level",
			l:    logging.LogLevel(99),
			want: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exp := tt.want
			act := tt.l.String()

			assert.Equalf(t, exp, act, "LogLevel.String() = %v, want %v", act, exp)
		})
	}
}

func TestLogLevel_Limit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    logging.LogLevel
		want logging.LogLevel
	}{
		{
			name: "debug level",
			l:    logging.LevelDebug,
			want: logging.LevelDebug,
		},
		{
			name: "info level",
			l:    logging.LevelInfo,
			want: logging.LevelInfo,
		},
		{
			name: "warning level",
			l:    logging.LevelWarn,
			want: logging.LevelWarn,
		},
		{
			name: "error level",
			l:    logging.LevelError,
			want: logging.LevelError,
		},
		{
			name: "fatal level",
			l:    logging.LevelFatal,
			want: logging.LevelError,
		},
		{
			name: "unknown level",
			l:    logging.LogLevel(99),
			want: logging.LevelError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exp := tt.want
			act := tt.l.Validate()

			assert.Equalf(t, exp, act, "LogLevel.Limit() = %v, want %v", act, exp)
		})
	}
}
