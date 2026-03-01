// Package server provides general functionality for running a server application.
package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/app"
	"github.com/mr-filatik/go-password-keeper/internal/platform/http/diagnostic"
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
	projectName = "go_password_keeper"
	appName     = "server"

	shutdownTimeout = 5 * time.Second
)

// // IServer - interface for all application servers.
// type IServer interface {
// 	// Starting the server.
// 	//
// 	// Implements the platform.IStarter interface.
// 	platform.IStarter

// 	// Correct server shutdown.
// 	//
// 	// Implements the platform.IShutdowner interface.
// 	platform.IShutdowner
// }

// Run starts the server application.
//
//nolint:funlen // Run() is the main function in which all components are initialized.
func Run() {
	logger, loggerErr := logging.NewZapSugarLoggerWithFields(
		logging.LevelWarn,
		os.Stdout,
		logging.FormatJSON,
		"project", projectName,
		"app", appName,
	)
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
		"build version", buildVersion,
		"build date", buildDate,
		"build commit", buildCommit,
	)

	// ===== Binding OS signals to context =====
	exitCtx, exitFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer exitFn()

	appConfig := config.Initialize()

	// ===== CREATING METRICS =====

	metricsProvider := metrics.CreateProvider("filatik", projectName, appName)

	metricsProvider.App.SetBuildInfo(metrics.AppBuildLabel{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	})

	metricsProvider.App.SetDeployInfo(metrics.AppDeployLabel{
		Number: "unknown",
	})

	// ===== DIAGNOSIC SERVER =====

	diagnosticServer := diagnostic.NewServer(diagnostic.ServerConfig{
		Address:         appConfig.DiagnosticAddress,
		MetricsProvider: metricsProvider,
	}, logger)

	diagnosticServerStartErr := diagnosticServer.Start(exitCtx)
	if diagnosticServerStartErr != nil {
		logger.Error("Starting diagnostic server error", diagnosticServerStartErr)

		return
	}

	// ===== CREATING SERVICES =====

	mainServer := http.NewServer(
		"main http server",
		http.ServerConfig{
			Address:         appConfig.Address,
			MetricsProvider: metricsProvider,
		}, logger)

	addServer := http.NewServer(
		"add http server",
		http.ServerConfig{
			Address:         ":31212",
			MetricsProvider: metricsProvider,
		}, logger)

	// cacher := redis.NewCacher(redis.CacherConfig{
	// 	ClientName:  "server",
	// 	Address:     "redis:6379",
	// 	DBNumber:    0,
	// 	Username:    "",
	// 	Password:    "",
	// 	ConnTimeout: 2 * time.Second,
	// }, logger)

	apper := app.New(logger, metricsProvider, diagnosticServer)

	apper.RegisterComponents(
		addServer,
		apper.WithParallelComponent(
			mainServer,
			apper.WithSequentialComponent(),
		),
	)

	startErr := apper.Start(exitCtx)
	if startErr != nil {
		logger.Error("Starting services error", startErr)
	}

	<-exitCtx.Done()

	// Нужно понять, как сделать так, чтобы при получении сигнала не закрывать сразу
	// А ждать это время
	// И указать время для каждого конкретного компонента
	shutdownCtx, cansel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cansel()

	shutdownErr := apper.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		logger.Error("Stoping services error", shutdownErr)
	}

	// // ===== STARTING SERVICES =====

	// starter := platform.NewStarter(
	// 	logger,
	// 	platform.NewSequentialStarter( /*cacher, */ mainServer),
	// 	metricsProvider,
	// )

	// startServicesErr := starter.Start(exitCtx)
	// if startServicesErr != nil {
	// 	logger.Error("Starting services error", startServicesErr)

	// 	exitFn()
	// 	//return
	// }

	// logger.Info("Application starting is successful")

	// // ===== Waiting for the stop signal =====
	// <-exitCtx.Done()

	// // ===== Start of server shutdown =====
	// shutdownCtx, cansel := context.WithTimeout(context.Background(), shutdownTimeout)
	// defer cansel()

	// // ===== STOPPING SERVICES =====

	// stopper := platform.NewStopper(
	// 	logger,
	// 	platform.NewSequentialStopper(logger, mainServer /*cacher, */),
	// 	metricsProvider,
	// 	diagnosticServer,
	// )

	// // для diagnosticServer shutdownCtx не нужен, нужно дополнительно
	// stopServicesErr := stopper.Shutdown(shutdownCtx)
	// if stopServicesErr != nil {
	// 	logger.Error("Stoppping services error", startServicesErr)
	// }
}
