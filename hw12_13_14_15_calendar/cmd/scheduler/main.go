package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	app "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewSchedulerConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Output)
	defer log.Close()

	scanInterval, err := time.ParseDuration(config.Schedule.Interval)
	if err != nil {
		log.Error(fmt.Sprintf("failed to parse interval value: %s", err.Error()))
		log.Close()
		os.Exit(1) //nolint:gocritic
	}

	retentionPeriod, err := time.ParseDuration(config.Schedule.RetentionPeriod)
	if err != nil {
		log.Error(fmt.Sprintf("failed to parse interval value: %s", err.Error()))
		log.Close()
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := storage.NewStorage(config.Storage)
	if err := storage.Connect(ctx); err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		log.Close()
		os.Exit(1)
	}
	defer storage.Close(ctx)

	queueConn := queue.NewConnection(config.Queue.QueueServerConf)
	if err := queueConn.Connect(); err != nil {
		log.Error(fmt.Sprintf("failed to connect to queue: %s", err.Error()))
		storage.Close(ctx)
		log.Close()
		os.Exit(1)
	}
	defer queueConn.Close()

	producer := queue.NewProducer(queueConn, config.Queue)
	if err := producer.Start(); err != nil {
		log.Error(fmt.Sprintf("failed to create exchange: %s", err.Error()))
		queueConn.Close()
		storage.Close(ctx)
		log.Close()
		os.Exit(1)
	}
	defer producer.Close()

	schedulerApp := app.New(log, storage, producer, scanInterval, retentionPeriod)

	log.Info("Scheduler is started")
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-ticker.C:
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				if err := schedulerApp.ProcessNotifications(ctx); err != nil {
					log.Error("failed to process notifications: %s", err)
				}
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				if err := schedulerApp.PurgeOldEvents(ctx); err != nil {
					log.Error("failed to purge events: %s", err)
				}
				wg.Done()
			}()

			wg.Wait()
		}
	}

	log.Info("Scheduler is stopped")
}
