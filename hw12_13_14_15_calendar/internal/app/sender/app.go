package senderapp

import (
	"context"
	"errors"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var errCannotBeEmpty = errors.New("cannot be empty")

type App struct {
	logger  common.Logger
	storage storage.Storage
}

func New(logger common.Logger, storage storage.Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) Notify(ctx context.Context, dto contracts.Notification) error {
	params := []any{
		"Event Id", dto.ID,
		"Time", dto.Time,
		"Title", dto.Title,
		"Email", dto.OwnerEmail,
	}
	if err := a.validateAttributes(dto); err != nil {
		return err
	}
	a.logger.Info("Notification sent", params...)

	err := a.storage.SetEventNotified(ctx, dto.ID)
	return err
}

func (a *App) validateAttributes(dto contracts.Notification) error {
	if dto.Title == "" {
		return customerrors.ValidationError{Field: "Title", Err: errCannotBeEmpty}
	}

	if dto.ID == "" {
		return customerrors.ValidationError{Field: "ID", Err: errCannotBeEmpty}
	}

	if dto.OwnerEmail == "" {
		return customerrors.ValidationError{Field: "Email", Err: errCannotBeEmpty}
	}

	return nil
}
