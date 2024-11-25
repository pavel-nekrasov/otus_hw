package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/http/handlers"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	host       string
	port       int
	logger     Logger
	app        Application
	httpServer *http.Server
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Application interface{}

func NewServer(host string, port int, logger Logger, app Application) *Server {
	return &Server{host: host, port: port, logger: logger, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	bindAddr := fmt.Sprintf("%v:%v", s.host, s.port)
	s.logger.Info(fmt.Sprintf("Starting on %v...", bindAddr))
	s.httpServer = &http.Server{Addr: bindAddr, ReadTimeout: time.Second * 10}

	helloHandler := handlers.NewHelloService(s.logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler.GetHello)

	muxWithLogging := middleware.NewLoggingMiddleware(s.logger, mux)

	s.httpServer.Handler = muxWithLogging

	s.httpServer.ListenAndServe()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// TODO
