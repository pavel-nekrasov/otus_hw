package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type (
	NewEventDto struct {
		id          string
		title       string
		start       string
		end         string
		description string
		owner       string
		notify      string
	}
)

func TestAppCreateEventSuccess(t *testing.T) {
	tests := []NewEventDto{
		{
			title:  "meeting 1",
			start:  "2024-11-25T10:00:00.000Z",
			end:    "2024-11-25T10:30:00.000Z",
			owner:  "user@example.com",
			notify: "",
		},
		{
			title:  "meeting 2",
			start:  "2024-11-25T12:00:00.000Z",
			end:    "2024-11-25T12:30:00.000Z",
			owner:  "user@example.com",
			notify: "",
		},
		{
			title:  "meeting 3",
			start:  "2024-11-25T10:00:00.000Z",
			end:    "2024-11-25T10:30:00.000Z",
			owner:  "user2@example.com",
			notify: "",
		},
		{
			title:  "meeting 4",
			start:  "2024-11-25T12:00:00.000Z",
			end:    "2024-11-25T12:30:00.000Z",
			owner:  "user2@example.com",
			notify: "",
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

			ev, err := app.CreateEvent(ctx, tt.title, tt.start, tt.end, tt.description, tt.owner, tt.notify)
			require.NoError(t, err)
			require.NotEmpty(t, ev.ID)
		})
	}
}

func TestAppCreateEventSimpleValidations(t *testing.T) {
	tests := []struct {
		add NewEventDto
		err error
	}{
		{
			add: NewEventDto{
				title:  "meeting 1",
				start:  "2024-11-25T10:00:00.000Z",
				end:    "2024-11-25T10:30:00",
				owner:  "user@example.com",
				notify: "",
			},
			err: customerrors.ValidationError{Field: "EndTime", Err: errors.New("wrong date/time format")},
		},
		{
			add: NewEventDto{
				title:  "meeting 2",
				start:  "2024-11-25T10:00:00",
				end:    "2024-11-25T10:30:00.000Z",
				owner:  "user@example.com",
				notify: "",
			},
			err: customerrors.ValidationError{Field: "StartTime", Err: errors.New("wrong date/time format")},
		},
		{
			add: NewEventDto{
				title:  "",
				start:  "2024-11-25T10:00:00.000Z",
				end:    "2024-11-25T10:30:00.000Z",
				owner:  "user@example.com",
				notify: "",
			},
			err: customerrors.ValidationError{Field: "Title", Err: errors.New("cannot be empty")},
		},
		{
			add: NewEventDto{
				title:  "Title 4",
				start:  "2024-11-25T10:00:00.000Z",
				end:    "2024-11-25T10:30:00.000Z",
				owner:  "",
				notify: "",
			},
			err: customerrors.ValidationError{Field: "OwnerEmail", Err: errors.New("cannot be empty")},
		},
		{
			add: NewEventDto{
				title:  "Title 5",
				start:  "2024-11-25T11:00:00.000Z",
				end:    "2024-11-25T10:30:00.000Z",
				owner:  "user@example.com",
				notify: "",
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

			_, err := app.CreateEvent(ctx, tt.add.title, tt.add.start, tt.add.end,
				tt.add.description, tt.add.owner, tt.add.notify)
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

	tests := []NewEventDto{
		{
			id:          "xxx2",
			title:       "meeting 2",
			start:       "2024-11-25T12:00:00.000Z",
			end:         "2024-11-25T12:30:00.000Z",
			owner:       "user@example.com",
			description: "new description2",
			notify:      "",
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
			_, err := app.CreateEvent(ctx, tt.title, tt.start, tt.end, tt.description, tt.owner, tt.notify)
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

	tests := []NewEventDto{
		{
			id:          "xxx",
			title:       "meeting 1",
			start:       "2024-11-25T10:00:00.000Z",
			end:         "2024-11-25T10:30:00.000Z",
			owner:       "user@example.com",
			description: "new description",
			notify:      "",
		},
		{
			id:          "xxx2",
			title:       "metting 2",
			start:       "2024-11-25T12:00:00.000Z",
			end:         "2024-11-25T12:30:00.000Z",
			owner:       "user@example.com",
			description: "new description2",
			notify:      "",
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
			err := app.UpdateEvent(ctx, tt.id, tt.title, tt.start, tt.end, tt.description, tt.owner, tt.notify)
			require.NoError(t, err)
			ev, err := app.GetEvent(ctx, tt.id)
			require.NoError(t, err)
			require.Equal(t, tt.description, ev.Description)
		})
	}
}

func TestAppUpdateEventNotFound(t *testing.T) {
	tests := []NewEventDto{
		{
			id:     "id1",
			title:  "meeting 1",
			start:  "2024-11-25T10:00:00.000Z",
			end:    "2024-11-25T10:30:00.000Z",
			owner:  "user@example.com",
			notify: "",
		},
		{
			id:     "id2",
			title:  "meeting 2",
			start:  "2024-11-25T12:00:00.000Z",
			end:    "2024-11-25T12:30:00.000Z",
			owner:  "user@example.com",
			notify: "",
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

			err := app.UpdateEvent(ctx, tt.id, tt.title, tt.start, tt.end, tt.description, tt.owner, tt.notify)
			require.Error(t, err)
		})
	}
}
