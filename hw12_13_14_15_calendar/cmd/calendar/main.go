package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewCalendarConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Output)
	defer log.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := storage.NewStorage(config.Storage)
	err := storage.Connect(ctx)
	if err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	calendar := app.New(log, storage)
	server := internalgrpc.NewServer(config.Endpoint.Host,
		config.Endpoint.GRPCPort,
		config.Endpoint.HTTPPort,
		log,
		calendar,
	)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop calendar server: " + err.Error())
		}
		log.Info("Shutted down")
	}()

	log.Info("Calendar is running...")
	if err := server.Start(ctx); err != nil {
		log.Error("failed to start calendar server: " + err.Error())
		cancel()
		os.Exit(1)
	}
	log.Info("Calendar is stopped")
}
