package storage

import (
	"context"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	sqlstorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type (
	Storage interface {
		AddEvent(ctx context.Context, event model.Event) error
		UpdateEvent(ctx context.Context, event model.Event) error
		GetEvent(ctx context.Context, eventID string) (model.Event, error)
		DeleteEvent(ctx context.Context, eventID string) error
		DeleteEventsOlderThan(ctx context.Context, time time.Time) error
		ListOwnerEventsForPeriod(ctx context.Context, ownerEmail string, startDate, endDate time.Time) ([]model.Event, error)
		ListEventsToBeNotified(ctx context.Context, startTime, endTime time.Time) ([]model.Event, error)
	}

	Control interface {
		Storage
		Connect(ctx context.Context) error
		Migrate(ctx context.Context, migrate string) (err error)
		Close(ctx context.Context) error
	}
)

func NewStorage(c config.StorageConf) Control {
	if c.Mode == config.StorageModePostgres {
		return sqlstorage.New(c.Host, c.Port, c.DBName, c.User, c.Password)
	}

	return memorystorage.New()
}
