// Package http contains a description of the HTTP server.
package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mr-filatik/go-password-keeper/docs/swagger/server" // Swagger docs registration in HTTP server.
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Server - describes the structure of an HTTP server.
type Server struct {
	router          *chi.Mux
	server          *http.Server
	metricsProvider *metrics.Provider
	logger          logging.Logger
	address         string
}

// ServerConfig - HTTP server configuration.
type ServerConfig struct {
	Address         string // Address
	MetricsProvider *metrics.Provider
}

const (
	timeoutIdle       = 5 * time.Second
	timeoutRead       = 5 * time.Second
	timeoutReadHeader = 5 * time.Second
	timeoutWrite      = 10 * time.Second
)

// NewServer - creates a new HTTP server instance.
func NewServer(conf ServerConfig, logger logging.Logger) *Server {
	logger.Info("Server creating...")

	tslNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	srvr := &Server{
		address:         conf.Address,
		metricsProvider: conf.MetricsProvider,
		logger:          logger,
		router:          chi.NewRouter(),
		server: &http.Server{
			Addr:                         conf.Address,
			BaseContext:                  nil,
			ConnContext:                  nil,
			ConnState:                    nil,
			DisableGeneralOptionsHandler: false,
			ErrorLog:                     nil,
			Handler:                      nil,
			IdleTimeout:                  timeoutIdle,
			MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
			ReadHeaderTimeout:            timeoutReadHeader,
			ReadTimeout:                  timeoutRead,
			TLSConfig:                    nil,
			TLSNextProto:                 tslNextProto,
			WriteTimeout:                 timeoutWrite,
			Protocols:                    nil,
			HTTP2:                        nil,
		},
	}

	srvr.registerMiddlewares()

	srvr.registerHandlers()

	logger.Info("Server create is successful")

	return srvr
}

// Start - starting the server.
//
// Implements the server.IServer interface.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(
		"Server starting...",
		"address", s.address,
	)

	s.server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error("Error in Server", err)
			} else {
				s.logger.Info("Server is closed")
			}
		}
	}()

	s.logger.Info("Server start is successful")

	return nil
}

// Shutdown gracefully terminates server.
//
// Implements the server.IServer interface.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Server shutdown starting...")

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	s.logger.Info("Server shutdown is successful")

	return nil
}

// Close - server shuts down.
//
// Implements the server.IServer interface.
func (s *Server) Close() error {
	s.logger.Info("Server close starting...")

	err := s.server.Close()
	if err != nil {
		return fmt.Errorf("server close: %w", err)
	}

	s.logger.Info("Server close is successful")

	return nil
}

func (s *Server) registerMiddlewares() {
	// 	r.Use(RecoverMiddleware)     // 1️⃣ ловит паники
	// r.Use(LoggingMiddleware)     // 2️⃣ логирует запрос/ответ
	// r.Use(MetricsMiddleware)     // 3️⃣ собирает метрики
	// r.Use(RequestIDMiddleware)   // 4️⃣ присваивает request id
	// r.Use(AuthMiddleware)        // 5️⃣ авторизация и т.д.

	// предлагается сделать один middleware для observability, который Recover, Logging и Metrics

	s.router.Use(
		middleware.Recover(s.logger),
		middleware.RequestID(),
		middleware.Logging(
			s.logger,
			//middleware.LoggingWithEnableBodyLogging(),
			middleware.LoggingWithRequestRoute(func(r *http.Request) string {
				return chi.RouteContext(r.Context()).RoutePattern()
			}),
		),
		middleware.Metrics(
			s.metricsProvider,
			middleware.MetricsWithRequestRoute(func(r *http.Request) string {
				return chi.RouteContext(r.Context()).RoutePattern()
			}),
		),
	)

	s.server.Handler = s.router
}

func (s *Server) registerHandlers() {
	s.router.Handle("/ping", http.HandlerFunc(s.ping))

	metrics.RegisterHandler(s.router)

	s.router.Handle("/swagger/*", httpSwagger.WrapHandler)

	s.server.Handler = s.router
}

const tempRandValue = 400

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	ok := s.validateRequestMethod(w, r.Method, http.MethodGet)
	if !ok {
		return
	}

	panic(errors.New("aaaa"))

	//start := time.Now()

	s.logger.Info("ping")

	w.WriteHeader(http.StatusOK)

	//nolint:gosec // temp code
	time.Sleep(time.Duration(rand.Int64N(tempRandValue)) * time.Millisecond)

	// s.metricsProvider.HTTP.IncRequestsCounter(metrics.HTTPRequestLabel{
	// 	Method:     r.Method,
	// 	Path:       "/ping",
	// 	StatusCode: http.StatusOK,
	// })
	// s.metricsProvider.HTTP.ObserveRequestDurationHistogram(metrics.HTTPRequestLabel{
	// 	Method:     r.Method,
	// 	Path:       "/ping",
	// 	StatusCode: http.StatusOK,
	// }, time.Since(start))

	_, err := w.Write([]byte("pong"))
	if err != nil {
		s.logger.Error("Internal server error (code 500)", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}
}

var errInvalidRequestMethod = errors.New("invalid request method")

func (s *Server) validateRequestMethod(w http.ResponseWriter, current string, needed string) bool {
	if current != needed {
		s.logger.Error(
			"Invalid request",
			errInvalidRequestMethod,
			"actual", current,
			"expected", needed,
		)

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		return false
	}

	return true
}
