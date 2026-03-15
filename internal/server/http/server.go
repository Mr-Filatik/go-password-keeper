// Package http contains a description of the HTTP server.
package http

//	@title					Password Keeper API
//	@version				1.0
//	@description.markdown	description
//	@termsOfService			https://example.com/terms

//	@contact.name	API Support
//	@contact.email	support@passwordkeeper.com

//	@license.name	Proprietary

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mr-filatik/go-password-keeper/docs/swagger/server" // Swagger docs registration in HTTP server.
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/handler"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Server - describes the structure of an HTTP server.
type Server struct {
	name string

	router          *chi.Mux
	server          *http.Server
	metricsProvider *metrics.Provider
	logger          logging.Logger
	address         string

	mu      sync.Mutex
	started bool
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
func NewServer(name string, conf ServerConfig, logger logging.Logger) *Server {
	logger.Info("Server creating...", "address", conf.Address)

	tslNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	srvr := &Server{
		name:            name,
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
		started: false,
		mu:      sync.Mutex{},
	}

	srvr.registerMiddlewares()

	srvr.registerHandlers()

	logger.Info("Server create is successful")

	return srvr
}

// GetName displays the name of the component.
//
// Implements the IComponent interface
// from "github.com/mr-filatik/go-password-keeper/internal/platform/app" package.
func (s *Server) GetName() string {
	return s.name
}

// Start - starting the server.
//
// Implements the server.IServer interface.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting...", nil, "component", s.GetName())

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		s.logger.Warn("Server already started", nil)

		return nil
	}

	s.logger.Info(
		"Starting HTTP server...",
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

	s.started = true

	return nil
}

// Shutdown gracefully terminates server.
//
// Implements the server.IServer interface.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Warn("Shutdowning...", nil, "component", s.GetName())

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		s.logger.Warn("Server not be started and not be stopped", nil)

		return nil
	}

	s.logger.Info("Server shutdown starting...")

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	s.logger.Info("Server shutdown is successful")

	s.started = false

	return nil
}

// Stop - server shuts down.
//
// Implements the server.IServer interface.
func (s *Server) Stop() error {
	s.logger.Warn("Stoping...", nil, "component", s.GetName())

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		s.logger.Warn("Server not be started and not be stopped", nil)

		return nil
	}

	s.logger.Info("Server stop starting...")

	err := s.server.Close()
	if err != nil {
		return fmt.Errorf("server stop: %w", err)
	}

	s.logger.Info("Server stop is successful")

	s.started = false

	return nil
}

func routeFromChiContext(r *http.Request) string {
	return chi.RouteContext(r.Context()).RoutePattern()
}

func (s *Server) registerMiddlewares() {
	s.router.Use(
		// middleware.Recover(s.logger), // сделать глобальный recover???
		middleware.InjectLogger(s.logger),
		middleware.RequestID(), // простой и не паникует
		// middleware.Limiter(...), // limiter: дешёво отстреливаем лишнее, защищает от DDoS / флудеров вообще.
		// middleware.LimiterUserID(...), // Пользовательский (по user_id) — уже после Auth, в защищённой группе.
		middleware.Recover(),
		// auth после recover, т.к. тут поход в базу, но ниже логинга, т.к.
		// он может положить user_id в контекст
		middleware.Logging(
			middleware.LoggingOpts{
				EnableRequestBodyLogging:  false,
				EnableResponseBodyLogging: false,
				RouteFn:                   routeFromChiContext,
			},
		),
		middleware.Metrics(
			s.metricsProvider,
			middleware.MetricsOpts{
				RouteFn: routeFromChiContext,
			},
		),
	)

	// Example:
	// s.router.Group(func(r chi.Router) {
	//     r.Use(middleware.Auth(authService, authOpts))

	//     r.Get("/me", getProfileHandler)
	//     r.Get("/orders", listOrdersHandler)
	//     r.Post("/orders", createOrderHandler)
	// })

	s.server.Handler = s.router
}

func (s *Server) registerHandlers() {
	s.router.Handle("/ping", http.HandlerFunc(s.ping))
	s.router.Post("/test", handler.Test())

	s.router.Handle("/swagger/*", httpSwagger.WrapHandler)

	s.server.Handler = s.router
}

const tempRandValue = 400

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	ok := s.validateRequestMethod(w, r.Method, http.MethodGet)
	if !ok {
		return
	}

	s.logger.Info("ping")

	w.WriteHeader(http.StatusOK)

	//nolint:gosec // temp code
	time.Sleep(time.Duration(rand.Int64N(tempRandValue)) * time.Millisecond)

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
