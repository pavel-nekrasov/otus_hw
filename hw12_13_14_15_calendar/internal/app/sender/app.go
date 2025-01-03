package senderapp

import (
	"context"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
)

type App struct {
	logger common.Logger
}

func New(logger common.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Notify(_ context.Context, dto contracts.Notification) error {
	params := []any{
		"Event Id", dto.ID,
		"Time", dto.Time,
		"Title", dto.Title,
		"Email", dto.OwnerEmail,
	}
	a.logger.Info("Notification sent", params...)
	return nil
}
