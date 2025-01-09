package schedulerapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
)

type App struct {
	logger          common.Logger
	storage         Storage
	publisher       Publisher
	scanInterval    time.Duration
	retentionPeriod time.Duration
	startTime       time.Time
}

type Storage interface {
	DeleteEventsOlderThan(ctx context.Context, time time.Time) error
	ListEventsToBeNotified(ctx context.Context, startTime, endTime time.Time) ([]model.Event, error)
}

type Publisher interface {
	Publish(data []byte) error
}

func New(logger common.Logger, storage Storage, publisher Publisher, scanInterval, retentionPeriod time.Duration) *App {
	return &App{
		logger:          logger,
		storage:         storage,
		publisher:       publisher,
		scanInterval:    scanInterval,
		retentionPeriod: retentionPeriod,
		startTime:       time.Now().Add(-scanInterval),
	}
}

func (a *App) ProcessNotifications(ctx context.Context) error {
	a.logger.Debug("Checking events to be notified...")

	events, err := a.storage.ListEventsToBeNotified(ctx, a.shiftBoundary(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to retrieve events to be notified: %w", err)
	}

	for _, ev := range events {
		dto := contracts.Notification{ID: ev.ID, Title: ev.Title, Time: ev.StartTime.Unix(), OwnerEmail: ev.OwnerEmail}
		a.logger.Debug(fmt.Sprintf("Publishing notification for event ID=%s", ev.ID))

		data, err := json.Marshal(dto)
		if err != nil {
			return fmt.Errorf("failed to serialize event: %w", err)
		}
		err = a.publisher.Publish(data)
		if err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	return nil
}

func (a *App) PurgeOldEvents(ctx context.Context) error {
	a.logger.Debug("Purging old events...")

	if err := a.storage.DeleteEventsOlderThan(ctx, time.Now().Add(-a.retentionPeriod)); err != nil {
		return fmt.Errorf("failed to delete old events: %w", err)
	}

	return nil
}

func (a *App) shiftBoundary() time.Time {
	defer func() {
		a.startTime = a.startTime.Add(a.scanInterval)
	}()

	return a.startTime
}
