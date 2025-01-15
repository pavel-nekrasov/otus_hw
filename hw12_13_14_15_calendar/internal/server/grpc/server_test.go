package internalgrpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	app "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/events"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	memorystorage "github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	storage := memorystorage.New()
	logger := logger.New("INFO", "stdout")
	app := app.New(logger, storage)
	service := events.NewService(logger, app)
	pb.RegisterEventsServer(s, service)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCRUDSuccess(t *testing.T) {
	cases := []*pb.TransientEvent{
		{
			Title:       "title 1",
			StartTime:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:     time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC).Unix(),
			Description: "description 1",
			OwnerEmail:  "user1@example.com",
			Notify:      "",
		},
		{
			Title:       "title 2",
			StartTime:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:     time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC).Unix(),
			Description: "description 2",
			OwnerEmail:  "user2@example.com",
			Notify:      "",
		},
	}
	ctx := context.Background()

	//nolint:staticcheck
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewEventsClient(conn)

	for _, c := range cases {
		resp, err := client.CreateEvent(ctx, &pb.NewEventRequest{Event: c})
		require.NoError(t, err)
		event := resp.GetEvent()
		require.NotNil(t, event)
		require.NotEmpty(t, event.Id)
		require.Equal(t, c.Title, event.Title)
		require.Equal(t, c.StartTime, event.StartTime)
		require.Equal(t, c.EndTime, event.EndTime)
		require.Equal(t, c.Description, event.Description)
		require.Equal(t, c.OwnerEmail, event.OwnerEmail)
		require.Equal(t, c.Notify, event.Notify)

		resp, err = client.GetEvent(ctx, &pb.EventIdRequest{Id: event.Id})
		require.NoError(t, err)
		event = resp.GetEvent()
		require.NotNil(t, event)
		require.NotEmpty(t, event.Id)
		require.Equal(t, c.Title, event.Title)
		require.Equal(t, c.StartTime, event.StartTime)
		require.Equal(t, c.EndTime, event.EndTime)
		require.Equal(t, c.Description, event.Description)
		require.Equal(t, c.OwnerEmail, event.OwnerEmail)
		require.Equal(t, c.Notify, event.Notify)

		c.Description = "updated description"
		resp, err = client.UpdateEvent(ctx, &pb.UpdateEventRequest{Id: event.Id, Event: c})
		require.NoError(t, err)
		event = resp.GetEvent()
		require.NotNil(t, event)
		require.NotEmpty(t, event.Id)
		require.Equal(t, c.Title, event.Title)
		require.Equal(t, c.StartTime, event.StartTime)
		require.Equal(t, c.EndTime, event.EndTime)
		require.Equal(t, c.Description, event.Description)
		require.Equal(t, c.OwnerEmail, event.OwnerEmail)
		require.Equal(t, c.Notify, event.Notify)

		_, err = client.DeleteEvent(ctx, &pb.EventIdRequest{Id: event.Id})
		require.NoError(t, err)
	}
}
