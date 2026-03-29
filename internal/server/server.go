// Package server provides general functionality for running a server application.
package server

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/app"
	"github.com/mr-filatik/go-password-keeper/internal/platform/caching/redis"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
	zaplog "github.com/mr-filatik/go-password-keeper/internal/platform/log/zap"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
	"github.com/mr-filatik/go-password-keeper/internal/server/config"
	"github.com/mr-filatik/go-password-keeper/internal/server/http"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/dto"
)

//nolint:gochecknoglobals // substitution of linker flags via -ldflags
var (
	buildVersion = "N/A" // Application build version.
	buildDate    = "N/A" // Application build date.
	buildCommit  = "N/A" // Application build commit.
)

const (
	namespace   = "filatik"
	projectName = "go_password_keeper"
	appName     = "server"

	shutdownTimeout   = 5 * time.Second
	connectionTimeout = 2 * time.Second
)

// Run starts the server application.
//
//nolint:funlen // Run() is the main function in which all components are initialized.
func Run() {
	logger, loggerErr := zaplog.NewZapSugarLogger(
		log.LevelDebug,
		log.WithGlobalFields(
			log.WithProjectField(projectName),
			log.WithAppField(appName),
		),
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

	// ===== Binding OS signals to context =====
	exitCtx, exitFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer exitFn()

	appConfig := config.Initialize()

	log.LogInfo(logger, "Application starting...",
		log.WithDataField(map[string]string{
			"version": buildVersion,
			"date":    buildDate,
			"commit":  buildCommit,
		}))

	log.LogInfo(logger, "TEST",
		log.WithDataFieldSanitised(dto.UserInfo{
			ID:       "asascscs",
			Email:    "123456789@12345.123",
			Password: "csnskcnsckn",
		}))

	// ===== CREATING METRICS =====

	metricsProvider := metrics.CreateProvider(namespace, projectName, appName)

	metricsProvider.App.SetBuildInfo(metrics.AppBuildLabel{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	})

	metricsProvider.App.SetDeployInfo(metrics.AppDeployLabel{
		Number: "unknown",
	})

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

	cacher := redis.NewCacher("redis cacher",
		redis.CacherConfig{
			ClientName:  "server",
			Address:     ":6379", // "redis:6379",
			DBNumber:    0,
			Username:    "",
			Password:    "",
			ConnTimeout: connectionTimeout,
		}, logger)

	// ===== APP RUN =====

	app := app.New(logger,
		app.WithStartStopMetrics(metricsProvider.App),
		app.WithDiagnosticServer(appConfig.DiagnosticAddress, metricsProvider),
	)

	app.RegisterComponents(
		cacher,
		addServer,
		app.WithParallelComponent(
			mainServer,
			app.WithSequentialComponent(),
		),
	)

	startErr := app.Start(exitCtx)
	if startErr != nil {
		logger.Error("Starting services error", startErr)

		// exitCtx - close
	}

	<-exitCtx.Done()

	// Нужно понять, как сделать так, чтобы при получении сигнала не закрывать сразу
	// А ждать это время
	// И указать время для каждого конкретного компонента
	shutdownCtx, cansel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cansel()

	shutdownErr := app.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		logger.Error("Stoping services error", shutdownErr)
	}
}
