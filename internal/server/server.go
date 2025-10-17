// Package server provides general functionality for running a server application.
package server

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
	"github.com/mr-filatik/go-password-keeper/internal/server/config"
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
//
//nolint:funlen // Run() is the main function in which all components are initialized.
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

	appConfig := config.Initialize()

	metricsProvider := metrics.CreateProvider("filatik_go_password_keeper", "server")

	httpServerConfig := http.ServerConfig{
		Address:         appConfig.Address,
		MetricsProvider: metricsProvider,
	}

	var mainServer IServer = http.NewServer(httpServerConfig, logger)

	mainServerStartErr := mainServer.Start(exitCtx)
	if mainServerStartErr != nil {
		logger.Error("Starting server error", mainServerStartErr)
	}

	logger.Info("Application starting is successful")

	repeatEverySecond(func(i int) {
		var word string
		var status string

		switch i {
		case 1:
			word = "one"
		case 2:
			word = "two"
		case 3:
			word = "three"
		case 4:
			word = "four"
		case 5:
			word = "five"
		default:
			word = "unknown"
		}

		metricsProvider.Experiment.IncDistributionsCounter(metrics.ExperimentDistributionLabel{
			ExperimentName: "first-experiment",
			BranchName:     word,
			Distributor:    "server",
		})

		switch i {
		case 1:
			word = "one"
			status = "success"
		case 2:
			word = "one"
			status = "failed"
		case 3:
			word = "two"
			status = "success"
		case 4:
			word = "two"
			status = "failed"
		case 5:
			word = "three"
			status = "success"
		default:
			word = "unknown"
		}

		metricsProvider.Experiment.IncDistributionsCounter(metrics.ExperimentDistributionLabel{
			ExperimentName: "second-experiment",
			BranchName:     word,
			Distributor:    "server",
		})

		metricsProvider.Experiment.IncExecutionsCounter(metrics.ExperimentExecutionLabel{
			ExperimentName: "second-experiment",
			BranchName:     word,
			Executor:       "server",
			Status:         status,
		})

		metricsProvider.Experiment.IncExecutionsCounter(metrics.ExperimentExecutionLabel{
			ExperimentName: "zero-experiment",
			BranchName:     word,
			Executor:       "server",
			Status:         status,
		})
	})

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

func repeatEverySecond(fn func(int)) {
	// инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		value := rand.Intn(5) + 1 // Intn(5) возвращает 0–4, поэтому +1
		fn(value)
	}
}
