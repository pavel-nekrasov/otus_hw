package main

import (
	"context"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type (
	StorageCtrl interface {
		Connect(ctx context.Context) error
		Migrate(ctx context.Context, migrate string) (err error)
		Close(ctx context.Context) error
	}
)

func NewStorage(c config.StorageConf) StorageCtrl {
	if c.Mode == config.StorageModePostgres {
		return sqlstorage.New(c.Host, c.Port, c.DBName, c.User, c.Password)
	}

	return memorystorage.New()
}
