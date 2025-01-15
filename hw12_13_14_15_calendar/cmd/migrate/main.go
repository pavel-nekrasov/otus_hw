package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
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
	logg := logger.New(config.Logger.Level, config.Logger.Output)
	defer logg.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := NewStorage(config.Storage)
	err := storage.Connect(ctx)
	if err != nil {
		logg.Error("failed to connect to storage: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	logg.Info("Appying migrations...")
	err = storage.Migrate(ctx, "migrations")
	if err != nil {
		logg.Error(fmt.Sprintf("Failed to apply migrations: %v", err))
		os.Exit(1)
	}
	logg.Info("Done")
}
