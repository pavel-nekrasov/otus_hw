package senderapp

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestAppNotifySuccess(t *testing.T) {
	tests := []contracts.Notification{
		{
			ID:         "id1",
			Title:      "meeting 1",
			Time:       time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
			OwnerEmail: "user@example.com",
		},
		{
			ID:         "id2",
			Title:      "meeting 2",
			Time:       time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			OwnerEmail: "user@example.com",
		},
	}

	logger := logger.New("INFO", "stdout")
	ctx := context.Background()
	app := New(logger)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := app.Notify(ctx, tt)
			require.NoError(t, err)
		})
	}
}

func TestAppNotifyError(t *testing.T) {
	tests := []struct {
		notification contracts.Notification
		err          error
	}{
		{
			notification: contracts.Notification{
				ID:         "",
				Title:      "meeting 1",
				Time:       time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				OwnerEmail: "user@example.com",
			},
			err: customerrors.ValidationError{Field: "ID", Err: errors.New("cannot be empty")},
		},
		{
			notification: contracts.Notification{
				ID:         "ID2",
				Title:      "",
				Time:       time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				OwnerEmail: "user@example.com",
			},
			err: customerrors.ValidationError{Field: "Title", Err: errors.New("cannot be empty")},
		},
		{
			notification: contracts.Notification{
				ID:         "ID3",
				Title:      "Title 3",
				Time:       time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				OwnerEmail: "",
			},
			err: customerrors.ValidationError{Field: "Email", Err: errors.New("cannot be empty")},
		},
	}

	logger := logger.New("INFO", "stdout")
	ctx := context.Background()
	app := New(logger)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := app.Notify(ctx, tt.notification)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.err.Error())
			var valErr customerrors.ValidationError
			require.ErrorAs(t, err, &valErr)
		})
	}
}
