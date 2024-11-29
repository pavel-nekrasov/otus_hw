package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/http"
)

var configFile, migrateonly string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&migrateonly, "migrateonly", "false", "apply migration and exit")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.New(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.Output)

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

	if migrateonly == "true" {
		logg.Info("Appying migrations...")
		err := storage.Migrate(ctx, "migrations")
		if err != nil {
			logg.Error(fmt.Sprintf("Failed to apply migrations: %v", err))
			os.Exit(1)
		}
		logg.Info("Done")
		return
	}

	// ev, err := storage.GetEvent(ctx, "event-1")
	// if err != nil {
	// 	logg.Error(fmt.Sprintf("Failed to get event:%v", err))
	// } else {
	// 	logg.Info("Event received",
	// 		"ID", ev.ID,
	// 		"Title", ev.Title,
	// 		"StartTime", ev.StartTime,
	// 		"EndTime", ev.EndTime,
	// 		"Desc", ev.Description,
	// 		"Notify", ev.NotifyBefore,
	// 		"Owner", ev.OwnerEmail)
	// }

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
