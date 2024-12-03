package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestAppCreateEventSuccess(t *testing.T) {
	tests := []contracts.Event{
		{
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			Title:        "meeting 3",
			StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user2@example.com",
			NotifyBefore: "",
		},
		{
			Title:        "meeting 4",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user2@example.com",
			NotifyBefore: "",
		},
	}

	logger := logger.New("INFO", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			ev, err := app.CreateEvent(ctx, tt)
			require.NoError(t, err)
			require.NotEmpty(t, ev.ID)
		})
	}
}

func TestAppCreateEventSimpleValidations(t *testing.T) {
	tests := []struct {
		add contracts.Event
		err error
	}{
		{
			add: contracts.Event{
				Title:        "",
				StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
			err: customerrors.ValidationError{Field: "Title", Err: errors.New("cannot be empty")},
		},
		{
			add: contracts.Event{
				Title:        "Title 4",
				StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
				OwnerEmail:   "",
				NotifyBefore: "",
			},
			err: customerrors.ValidationError{Field: "OwnerEmail", Err: errors.New("cannot be empty")},
		},
		{
			add: contracts.Event{
				Title:        "Title 5",
				StartTime:    time.Date(2024, 11, 25, 11, 0, 0, 0, time.UTC).Unix(),
				EndTime:      time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
			err: customerrors.ValidationError{Field: "StartTime", Err: errors.New("must be less than EndTime")},
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := app.CreateEvent(ctx, tt.add)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.err.Error())
			var valErr customerrors.ValidationError
			require.ErrorAs(t, err, &valErr)
		})
	}
}

func TestAppCreateEventTimeBusy(t *testing.T) {
	data := []storage.Event{
		{
			ID:           "xxx",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
	}

	tests := []contracts.Event{
		{
			ID:           "xxx2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			Description:  "new description2",
			NotifyBefore: "",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)

	for _, ev := range data {
		err := storage.AddEvent(ctx, ev)
		require.NoError(t, err)
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			_, err := app.CreateEvent(ctx, tt)
			require.Error(t, err)
			require.ErrorContains(t, err, "StartTime|EndTime: another meeting exists for that period")
			var valErr customerrors.ValidationError
			require.ErrorAs(t, err, &valErr)
		})
	}
}

func TestAppGetEventSuccess(t *testing.T) {
	tests := []storage.Event{
		{
			ID:           "xxx",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			ID:           "xxx2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := storage.AddEvent(ctx, tt)
			require.NoError(t, err)
			ev, err := app.GetEvent(ctx, tt.ID)
			require.NoError(t, err)
			require.Equal(t, tt.ID, ev.ID)
			require.Equal(t, tt.Title, ev.Title)
			require.Equal(t, tt.StartTime, ev.StartTime)
			require.Equal(t, tt.EndTime, ev.EndTime)
			require.Equal(t, tt.Description, ev.Description)
			require.Equal(t, tt.NotifyBefore, ev.NotifyBefore)
			require.Equal(t, tt.OwnerEmail, ev.OwnerEmail)
		})
	}
}

func TestAppGetEventNotFound(t *testing.T) {
	tests := []struct {
		id string
	}{
		{
			id: "xxx1",
		},
		{
			id: "xxx2",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			_, err := app.GetEvent(ctx, tt.id)
			require.Error(t, err)
			var notFoundErr customerrors.NotFound
			require.ErrorAs(t, err, &notFoundErr)
		})
	}
}

func TestAppDeleteEventSuccess(t *testing.T) {
	tests := []storage.Event{
		{
			ID:           "xxx",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			ID:           "xxx2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := storage.AddEvent(ctx, tt)
			require.NoError(t, err)
			err = app.DeleteEvent(ctx, tt.ID)
			require.NoError(t, err)
		})
	}
}

func TestAppDeleteEventNotFound(t *testing.T) {
	tests := []struct {
		id string
	}{
		{
			id: "xxx1",
		},
		{
			id: "xxx2",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := app.DeleteEvent(ctx, tt.id)
			require.Error(t, err)
			var notFoundErr customerrors.NotFound
			require.ErrorAs(t, err, &notFoundErr)
		})
	}
}

func TestAppUpdateEventSuccess(t *testing.T) {
	data := []storage.Event{
		{
			ID:           "xxx",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			ID:           "xxx2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
	}

	tests := []contracts.Event{
		{
			ID:           "xxx",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			Description:  "new description",
			NotifyBefore: "",
		},
		{
			ID:           "xxx2",
			Title:        "metting 2",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			Description:  "new description2",
			NotifyBefore: "",
		},
	}

	logger := logger.New("DEBUG", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)

	for _, ev := range data {
		err := storage.AddEvent(ctx, ev)
		require.NoError(t, err)
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			_, err := app.UpdateEvent(ctx, tt)
			require.NoError(t, err)
			ev, err := app.GetEvent(ctx, tt.ID)
			require.NoError(t, err)
			require.Equal(t, tt.Description, ev.Description)
		})
	}
}

func TestAppUpdateEventNotFound(t *testing.T) {
	tests := []contracts.Event{
		{
			ID:           "id1",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 11, 25, 10, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 10, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
		{
			ID:           "id2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:      time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:   "user@example.com",
			NotifyBefore: "",
		},
	}

	logger := logger.New("INFO", "stdout")
	ctx := context.Background()
	storage := memorystorage.New()
	app := New(logger, storage)
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := app.UpdateEvent(ctx, tt)
			require.Error(t, err)
		})
	}
}
