package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/events"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/middleware"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	httpmiddleware "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/http/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	host     string
	grpcPort int
	httpPort int
	logger   Logger
	app      events.Application
	server   *grpc.Server
	gwServer *http.Server
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

func NewServer(host string, grpcPort int, httpPort int, logger Logger, app events.Application) *Server {
	return &Server{
		host:     host,
		grpcPort: grpcPort,
		httpPort: httpPort,
		logger:   logger,
		app:      app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	grpcBindAddr := fmt.Sprintf("%v:%v", s.host, s.grpcPort)
	s.logger.Info(fmt.Sprintf("Starting GRPC on %v...", grpcBindAddr))
	lsn, err := net.Listen("tcp", grpcBindAddr)
	if err != nil {
		return err
	}

	httpBindAddr := fmt.Sprintf("%v:%v", s.host, s.httpPort)
	s.logger.Info(fmt.Sprintf("Starting HTTP on %v...", httpBindAddr))

	gwClient, err := grpc.NewClient(grpcBindAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterEventsHandler(ctx, gwMux, gwClient)
	if err != nil {
		return err
	}

	gwMuxWithLogging := httpmiddleware.NewLoggingMiddleware(s.logger, gwMux)
	s.gwServer = &http.Server{Addr: httpBindAddr, Handler: gwMuxWithLogging, ReadTimeout: time.Second * 10}
	go func() {
		if err := s.gwServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(middleware.InterceptorLogger(s.logger)),
		),
	)
	pb.RegisterEventsServer(s.server, events.NewService(s.logger, s.app))
	reflection.Register(s.server)

	go func() {
		if err := s.server.Serve(lsn); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping adapters...")

	defer s.server.Stop()
	return s.gwServer.Shutdown(ctx)
}
