// Package adapter provides adapter functionality for the redis package.
package adapter

import (
	"context"
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
)

// LoggerAdapter - adapter for converting the redis logger to an application logger.
type LoggerAdapter struct {
	logging.Logger
}

// NewLoggerAdapter creates a new *LoggerAdapter instance.
//
// Parameters:
//   - logger logging.Logger: logger.
func NewLoggerAdapter(logger logging.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		Logger: logger,
	}
}

const formatMaintNotificationsDisabled = "auto mode fallback: maintnotifications disabled due to handshake error: %v"

// Printf prints a message from the redis library via the logger.
//
// Parameters:
//   - _ context.Context: context;
//   - format string: format string;
//   - vals ...interface{}: additional values.
//
// Implements the internal.Logging interface from github.com/redis/go-redis/v9.
func (l *LoggerAdapter) Printf(_ context.Context, format string, vals ...interface{}) {
	for index := range vals {
		switch val := vals[index].(type) {
		case error:
			//nolint:err113 // the format comes from another library
			err := fmt.Errorf(format, val)

			if format == formatMaintNotificationsDisabled {
				l.Warn("Issue from package github.com/redis/go-redis/v9", err)
			} else {
				l.Error("Issue from package github.com/redis/go-redis/v9", err)
			}

		case string:
			l.Info(fmt.Sprintf("Issue from package github.com/redis/go-redis/v9: %v", fmt.Sprintf(format, vals...)))
		}
	}
}

// Error examples:
//
// format:
//
// error:
//

// Warning examples:
//
// format:
// auto mode fallback: maintnotifications disabled due to handshake error: %v
// error:
// ERR unknown subcommand 'maint_notifications'. Try CLIENT HELP.
