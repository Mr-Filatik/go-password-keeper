// Package server provides general functionality for running a server application.
package server

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/server/http"
)

//nolint:gochecknoglobals // substitution of linker flags via -ldflags
var (
	buildVersion = "N/A" // Application build version.
	buildDate    = "N/A" // Application build date.
	buildCommit  = "N/A" // Application build commit.
)

const (
	shutdownTimeout = 5 * time.Second
)

// IServer - interface for all application servers.
type IServer interface {
	// Starting the server.
	//
	// Implements the platform.IStarter interface.
	platform.IStarter

	// Correct server shutdown.
	//
	// Implements the platform.IShutdowner interface.
	platform.IShutdowner
}

// Run starts the server application.
func Run() {
	logger, loggerErr := logging.NewZapSugarLogger(logging.LevelInfo, os.Stdout, logging.FormatJSON)
	if loggerErr != nil {
		panic(loggerErr)
	}

	defer func() {
		loggerErr := logger.Close()
		if loggerErr != nil {
			panic(loggerErr)
		}
	}()

	logger.Info("Application starting...",
		"Build Version", buildVersion,
		"Build Date", buildDate,
		"Build Commit", buildCommit,
	)

	// ===== Binding OS signals to context =====
	exitCtx, exitFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer exitFn()

	httpServerConfig := http.ServerConfig{
		Address: "localhost:8080",
	}

	var mainServer IServer = http.NewServer(httpServerConfig, logger)

	mainServerStartErr := mainServer.Start(exitCtx)
	if mainServerStartErr != nil {
		logger.Error("Starting server error", mainServerStartErr)
	}

	logger.Info("Application starting is successful")

	// ===== Waiting for the stop signal =====
	<-exitCtx.Done()

	// ===== Start of server shutdown =====
	logger.Info("Application shutdown starting...")

	shutdownCtx, cansel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cansel()

	shutdownErr := mainServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		logger.Error("Shutdown server error", shutdownErr)

		if errors.Is(shutdownErr, context.DeadlineExceeded) {
			logger.Warn("Shutdown context deadline exceeded, forcing close...", nil)
		}

		closeErr := mainServer.Close()
		if closeErr != nil {
			logger.Error("Close server error", closeErr)
		}
	}

	logger.Info("Application shutdown is successful")
}
