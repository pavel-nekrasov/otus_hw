package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/queue"
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

	config := config.NewSenderConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Output)
	defer log.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	senderApp := app.New(log)

	queueConn := queue.NewConnection(config.Queue.QueueServerConf)
	if err := queueConn.Connect(); err != nil {
		log.Error(fmt.Sprintf("failed to connect to queue: %s", err.Error()))
		log.Close()
		os.Exit(1) //nolint:gocritic
	}
	defer queueConn.Close()

	handler := queue.NewHandler(func(data []byte) error {
		var notification contracts.Notification
		if err := json.Unmarshal(data, &notification); err != nil {
			log.Error(fmt.Sprintf("failed to parse notification: %s", err))
			return err
		}
		if err := senderApp.Notify(ctx, notification); err != nil {
			log.Error(fmt.Sprintf("failed to send notification: %s", err))
			return err
		}
		return nil
	})

	consumer := queue.NewConsumer(queueConn, config.Queue, handler)
	defer consumer.Close()

	log.Info("Sender app started")
	if err := consumer.Start(ctx); err != nil {
		log.Error(fmt.Sprintf("failed to connect to queue: %s", err.Error()))
		consumer.Close()
		queueConn.Close()
		log.Close()
		os.Exit(1)
	}

	log.Info("Sender app stopped")
}
