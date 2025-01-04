package schedulerapp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/stretchr/testify/require"
)

type Mock struct {
	Data       []model.Event
	deleteCnt  int
	publishCnt int
}

func (t *Mock) ListEventsToBeNotified(_ context.Context, _, _ time.Time) ([]model.Event, error) {
	return t.Data, nil
}

func (t *Mock) DeleteEventsOlderThan(_ context.Context, _ time.Time) error {
	t.deleteCnt++
	return nil
}

func (t *Mock) Publish(_ []byte) error {
	t.publishCnt++
	return nil
}

func TestAppProcessSuccess(t *testing.T) {
	tests := []struct {
		Data               []model.Event
		expectedDeleteCnt  int
		expectedPublishCnt int
	}{
		{
			Data: []model.Event{
				{ID: "id1"},
				{ID: "id2"},
			},
			expectedDeleteCnt:  1,
			expectedPublishCnt: 2,
		},
		{
			Data: []model.Event{
				{ID: "id1"},
			},
			expectedDeleteCnt:  1,
			expectedPublishCnt: 1,
		},
		{
			Data:               []model.Event{},
			expectedDeleteCnt:  1,
			expectedPublishCnt: 0,
		},
	}

	logger := logger.New("INFO", "stdout")
	ctx := context.Background()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			mock := &Mock{}
			app := New(logger, mock, mock, 1*time.Second, 1*time.Second)
			mock.Data = tt.Data
			app.Process(ctx)
			require.Equal(t, tt.expectedDeleteCnt, mock.deleteCnt)
			require.Equal(t, tt.expectedPublishCnt, mock.publishCnt)
		})
	}
}
