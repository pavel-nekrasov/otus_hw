package app

import (
	"context"
	"errors"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	uuid "github.com/satori/go.uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Storage interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, eventID string) (storage.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error
	ListEventsForPeriod(ctx context.Context, ownerEmail string, startDate, endDate time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

var (
	errWrongDateFormat = errors.New("wrong date/time format")
	errCannotBeEmpty   = errors.New("cannot be empty")
	errWrongPeriod     = errors.New("must be less than EndTime")
	errNotSameDay      = errors.New("must be on the same day with EndTime")
	errPeriodIsBusy    = errors.New("another meeting exists for that period")
)

func (a *App) CreateEvent(ctx context.Context, dto contracts.Event) (storage.Event, error) {
	event, err := a.validateAttributes(ctx, dto)
	if err != nil {
		return storage.Event{}, err
	}
	event.ID = uuid.NewV4().String()
	err = a.storage.AddEvent(ctx, event)
	if err != nil {
		return storage.Event{}, err
	}
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, dto contracts.Event) error {
	event, err := a.validateAttributes(ctx, dto)
	if err != nil {
		return err
	}
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	return a.storage.GetEvent(ctx, eventID)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) ListEventsForDate(ctx context.Context, ownerEmail string, date string) ([]storage.Event, error) {
	dt, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, customerrors.ParamError{Param: "date", Err: errWrongDateFormat}
	}
	return a.storage.ListEventsForPeriod(ctx, ownerEmail, dt, dt)
}

func (a *App) validateAttributes(ctx context.Context, dto contracts.Event) (storage.Event, error) {
	if dto.Title == "" {
		return storage.Event{}, customerrors.ValidationError{Field: "Title", Err: errCannotBeEmpty}
	}

	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return storage.Event{}, customerrors.ValidationError{Field: "StartTime", Err: errWrongDateFormat}
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return storage.Event{}, customerrors.ValidationError{Field: "EndTime", Err: errWrongDateFormat}
	}

	if startTime.After(endTime) {
		return storage.Event{}, customerrors.ValidationError{Field: "StartTime", Err: errWrongPeriod}
	}

	if startTime.Day() != endTime.Day() || startTime.Month() != endTime.Month() || startTime.Year() != endTime.Year() {
		return storage.Event{}, customerrors.ValidationError{Field: "StartTime", Err: errNotSameDay}
	}

	if dto.OwnerEmail == "" {
		return storage.Event{}, customerrors.ValidationError{Field: "OwnerEmail", Err: errCannotBeEmpty}
	}

	if dto.NotifyBefore != "" {
		if _, err := time.Parse(time.RFC3339, dto.NotifyBefore); err != nil {
			return storage.Event{}, customerrors.ValidationError{Field: "Notify", Err: errWrongDateFormat}
		}
	}

	existingEvents, err := a.storage.ListEventsForPeriod(ctx, dto.OwnerEmail, startTime, endTime)
	if err != nil {
		return storage.Event{}, err
	}
	for _, ev := range existingEvents {
		if (startTime.Equal(ev.StartTime) || startTime.After(ev.StartTime)) &&
			(endTime.Equal(ev.EndTime) || endTime.Before(ev.EndTime)) {
			return storage.Event{}, customerrors.ValidationError{Field: "StartTime|EndTime", Err: errPeriodIsBusy}
		}
	}

	return storage.Event{
		ID:           dto.ID,
		Title:        dto.Title,
		StartTime:    startTime,
		EndTime:      endTime,
		Description:  dto.Description,
		OwnerEmail:   dto.OwnerEmail,
		NotifyBefore: dto.NotifyBefore,
	}, nil
}
