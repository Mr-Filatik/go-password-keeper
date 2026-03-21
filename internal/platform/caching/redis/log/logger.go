// Package redislog contains an implementation of the internal.Logging interface
// from the github.com/redis/go-redis/v9 package for integration with the logger
// used in the application.
package redislog

import (
	"context"
	"errors"
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
)

// LoggerAdapter provides logging methods adapted to the internal.Logging interface
// from the github.com/redis/go-redis/v9 package.
type LoggerAdapter struct {
	log.ILogger
}

// NewLoggerAdapter creates a new adapter instance.
func NewLoggerAdapter(logger log.ILogger) *LoggerAdapter {
	return &LoggerAdapter{
		ILogger: logger,
	}
}

const formatMaintNotificationsDisabled = "auto mode fallback: maintnotifications disabled due to handshake error: %v"

// Printf implements the internal.Logging interface from the github.com/redis/go-redis/v9 package,
// redirecting logs to the logger used by the application.
func (l *LoggerAdapter) Printf(_ context.Context, format string, vals ...interface{}) {
	var (
		errs   []error
		noerrs []any
	)

	for _, val := range vals {
		err, ok := val.(error)
		if ok {
			errs = append(errs, err)
		} else {
			noerrs = append(noerrs, val)
		}
	}

	msg := "Issue from package github.com/redis/go-redis/v9"

	if len(noerrs) > 0 {
		redisMsg := fmt.Sprintf(format, noerrs...)
		msg = fmt.Sprintf("Issue from package github.com/redis/go-redis/v9: %v", redisMsg)
	}

	if len(errs) == 0 {
		log.LogInfo(l, msg)

		return
	}

	//nolint:err113 // the format comes from external package
	err := fmt.Errorf("redis package error: %w", fmt.Errorf(format, errors.Join(errs...)))

	if format == formatMaintNotificationsDisabled {
		log.LogWarn(l, msg, err)
	} else {
		log.LogError(l, msg, err)
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
