package log_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    log.LogLevel
		want string
	}{
		{
			name: "debug level",
			l:    log.LevelDebug,
			want: "debug",
		},
		{
			name: "info level",
			l:    log.LevelInfo,
			want: "info",
		},
		{
			name: "warning level",
			l:    log.LevelWarn,
			want: "warn",
		},
		{
			name: "error level",
			l:    log.LevelError,
			want: "error",
		},
		{
			name: "fatal level",
			l:    log.LevelFatal,
			want: "fatal",
		},
		{
			name: "unknown level",
			l:    log.LogLevel(99),
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
		l    log.LogLevel
		want log.LogLevel
	}{
		{
			name: "debug level",
			l:    log.LevelDebug,
			want: log.LevelDebug,
		},
		{
			name: "info level",
			l:    log.LevelInfo,
			want: log.LevelInfo,
		},
		{
			name: "warning level",
			l:    log.LevelWarn,
			want: log.LevelWarn,
		},
		{
			name: "error level",
			l:    log.LevelError,
			want: log.LevelError,
		},
		{
			name: "fatal level",
			l:    log.LevelFatal,
			want: log.LevelError,
		},
		{
			name: "unknown level",
			l:    log.LogLevel(99),
			want: log.LevelError,
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
