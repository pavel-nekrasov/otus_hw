package main

import (
	"context"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type (
	Storage interface {
		AddEvent(ctx context.Context, event storage.Event) error
		UpdateEvent(ctx context.Context, event storage.Event) error
		GetEvent(ctx context.Context, eventID string) (storage.Event, error)
		DeleteEvent(ctx context.Context, eventID string) error
		ListEventsForPeriod(ctx context.Context, ownerEmail string, startDate, endDate time.Time) ([]storage.Event, error)
	}

	StorageCtrl interface {
		Storage
		Connect(ctx context.Context) error
		Migrate(ctx context.Context, migrate string) (err error)
		Close(ctx context.Context) error
	}
)

func NewStorage(config config.StorageConf) StorageCtrl {
	if config.Mode == "postgre" {
		return sqlstorage.New(config.Host, config.Port, config.DBName, config.User, config.Password)
	}

	return memorystorage.New()
}
