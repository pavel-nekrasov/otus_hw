package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageAddSuccess(t *testing.T) {
	tests := []struct {
		add         storage.Event
		expectedErr error
	}{
		{
			add: storage.Event{
				ID:           "id1",
				Title:        "meeting 1",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
		},
		{
			add: storage.Event{
				ID:           "id2",
				Title:        "meeting 2",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T12:00:00Z00:00",
			},
		},
		{
			add: storage.Event{
				ID:           "id3",
				Title:        "meeting 3",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T11:00:00Z00:00",
			},
		},
		{
			add: storage.Event{
				ID:           "id4",
				Title:        "meeting 4",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T12:00:00Z00:00",
			},
		},
	}

	ctx := context.Background()
	storage := New()
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := storage.AddEvent(ctx, tt.add)
			require.NoError(t, err)
			ev, err := storage.GetEvent(ctx, tt.add.ID)
			require.NoError(t, err)
			require.Equal(t, tt.add.Title, ev.Title)
			require.Equal(t, tt.add.StartTime, ev.StartTime)
			require.Equal(t, tt.add.EndTime, ev.EndTime)
			require.Equal(t, tt.add.Description, ev.Description)
			require.Equal(t, tt.add.OwnerEmail, ev.OwnerEmail)
			require.Equal(t, tt.add.NotifyBefore, ev.NotifyBefore)
		})
	}
}

func TestStorageGetSuccess(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		id       string
		expected storage.Event
	}{
		{
			id: "xxx",
			expected: storage.Event{
				ID:           "xxx",
				Title:        "meeting 1",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
		},
		{
			id: "xxx2",
			expected: storage.Event{
				ID:           "xxx2",
				Title:        "meeting 2",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T12:00:00Z00:00",
			},
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			ev, err := storage.GetEvent(ctx, tt.id)
			require.NoError(t, err)
			require.Equal(t, tt.expected.Title, ev.Title)
			require.Equal(t, tt.expected.StartTime, ev.StartTime)
			require.Equal(t, tt.expected.EndTime, ev.EndTime)
			require.Equal(t, tt.expected.Description, ev.Description)
			require.Equal(t, tt.expected.OwnerEmail, ev.OwnerEmail)
			require.Equal(t, tt.expected.NotifyBefore, ev.NotifyBefore)
		})
	}
}

func TestStorageGetError(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		id  string
		err error
	}{
		{
			id:  "xxx1",
			err: customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", "xxx1")},
		},
		{
			id:  "xxx3",
			err: customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", "xxx1")},
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := storage.GetEvent(ctx, tt.id)
			var custErr customerrors.NotFound
			require.ErrorAs(t, err, &custErr)
		})
	}
}

func TestStorageDeleteSuccess(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		id string
	}{
		{
			id: "xxx",
		},
		{
			id: "xxx2",
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := storage.DeleteEvent(ctx, tt.id)
			require.NoError(t, err)

			_, err = storage.GetEvent(ctx, tt.id)
			var custErr customerrors.NotFound
			require.ErrorAs(t, err, &custErr)
		})
	}
}

func TestStorageDeleteError(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		id  string
		err error
	}{
		{
			id:  "xxx1",
			err: customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", "xxx1")},
		},
		{
			id:  "xxx3",
			err: customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", "xxx1")},
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := storage.DeleteEvent(ctx, tt.id)
			var custErr customerrors.NotFound
			require.ErrorAs(t, err, &custErr)
		})
	}
}

func TestStorageUpdateSuccess(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		update storage.Event
	}{
		{
			update: storage.Event{
				ID:           "xxx",
				Title:        "meeting 1",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
		},
		{
			update: storage.Event{
				ID:           "xxx2",
				Title:        "meeting 2",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T12:00:00Z00:00",
			},
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := storage.UpdateEvent(ctx, tt.update)
			require.NoError(t, err)

			ev, err := storage.GetEvent(ctx, tt.update.ID)
			require.NoError(t, err)
			require.Equal(t, tt.update.Title, ev.Title)
			require.Equal(t, tt.update.StartTime, ev.StartTime)
			require.Equal(t, tt.update.EndTime, ev.EndTime)
			require.Equal(t, tt.update.Description, ev.Description)
			require.Equal(t, tt.update.OwnerEmail, ev.OwnerEmail)
			require.Equal(t, tt.update.NotifyBefore, ev.NotifyBefore)
		})
	}
}

func TestStorageUpdateError(t *testing.T) {
	prerequisites := []storage.Event{
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
	tests := []struct {
		update storage.Event
	}{
		{
			update: storage.Event{
				ID:           "xxx1",
				Title:        "meeting 1",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "",
			},
		},
		{
			update: storage.Event{
				ID:           "xxx3",
				Title:        "meeting 2",
				StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				OwnerEmail:   "user@example.com",
				NotifyBefore: "2024-01-01T12:00:00Z00:00",
			},
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := storage.UpdateEvent(ctx, tt.update)
			var custErr customerrors.NotFound
			require.ErrorAs(t, err, &custErr)
		})
	}
}

func TestListEventsForRang(t *testing.T) {
	prerequisites := []storage.Event{
		{
			ID:           "xxx1.1",
			Title:        "meeting 1",
			StartTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user1@example.com",
			NotifyBefore: "",
		},
		{
			ID:           "xxx1.2",
			Title:        "meeting 2",
			StartTime:    time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC),
			OwnerEmail:   "user1@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
		{
			ID:           "xxx1.3",
			Title:        "meeting 3",
			StartTime:    time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user1@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
		{
			ID:           "xxx3",
			Title:        "meeting 4",
			StartTime:    time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC),
			EndTime:      time.Date(2024, 1, 3, 12, 30, 0, 0, time.UTC),
			OwnerEmail:   "user2@example.com",
			NotifyBefore: "2024-01-01T12:00:00Z00:00",
		},
	}
	tests := []struct {
		user        string
		start       time.Time
		end         time.Time
		expectedCnt int
	}{
		{
			user:        "user1@example.com",
			start:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedCnt: 2,
		},
		{
			user:        "user1@example.com",
			start:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			expectedCnt: 1,
		},
		{
			user:        "user2@example.com",
			start:       time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			expectedCnt: 1,
		},
		{
			user:        "user3@example.com",
			start:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedCnt: 0,
		},
	}

	ctx := context.Background()
	storage := New()
	for _, e := range prerequisites {
		err := storage.AddEvent(ctx, e)
		require.NoError(t, err)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result, err := storage.ListEventsForPeriod(ctx, tt.user, tt.start, tt.end)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCnt, len(result))
		})
	}
}
