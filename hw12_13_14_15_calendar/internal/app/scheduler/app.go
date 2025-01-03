package schedulerapp

import (
	"context"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

type App struct {
	logger          common.Logger
	storage         Storage
	retentionPeriod time.Duration
}

type Storage interface {
	DeleteEvent(ctx context.Context, eventID string) error
	ListEventsForPeriod(ctx context.Context, ownerEmail string, startDate, endDate time.Time) ([]model.Event, error)
}

func New(logger common.Logger, storage Storage, retentionPeriod time.Duration) *App {
	return &App{logger: logger, storage: storage, retentionPeriod: retentionPeriod}
}

func (a *App) Process(_ context.Context) {
	a.logger.Info("Processing")
}
