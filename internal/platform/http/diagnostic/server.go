// Package diagnostic describes a diagnostic server for applications.
package diagnostic

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// Server - describes the structure of an HTTP server.
type Server struct {
	router          *chi.Mux
	server          *http.Server
	metricsProvider *metrics.Provider
	logger          log.ILogger
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
func NewServer(conf ServerConfig, logger log.ILogger) *Server {
	tslNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	srvr := &Server{
		address:         conf.Address,
		metricsProvider: conf.MetricsProvider,
		logger:          logger.With(log.WithComponentField("diagnostic server")),
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

	// srvr.registerMiddlewares()

	srvr.registerHandlers()

	return srvr
}

// Start - starting the server.
//
// Implements the server.IServer interface.
func (s *Server) Start(ctx context.Context) error {
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

	log.LogInfo(s.logger, "Server start is successful",
		log.WithDataField(map[string]string{
			"address": s.address,
		}))

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

func (s *Server) registerHandlers() {
	// metrics.RegisterHandler(s.router)
	metrics.RegisterHandlerWithScrapeCount(s.router)

	// // --- Health / ready ---
	// mux.HandleFunc("/healthz", s.handleHealth) // liveness
	// mux.HandleFunc("/readyz", s.handleReady)   // readiness

	// // --- Pprof ---
	// mux.HandleFunc("/debug/pprof/", pprof.Index)
	// mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.server.Handler = s.router
}
