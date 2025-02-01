package intgration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalendarIntegrationSuite struct {
	suite.Suite
	logger         *logger.Logger
	calendarConfig config.CalendarConfig
	scheduleConfig config.SchedulerConfig
	storage        storage.Control
}

const (
	CalendarHost = "cal_server_test"
)

func TestAppIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CalendarIntegrationSuite))
}

func (s *CalendarIntegrationSuite) SetupSuite() {
	s.calendarConfig = config.NewCalendarConfig("/app/configs/calendar_config.toml")
	s.scheduleConfig = config.NewSchedulerConfig("/app/configs/scheduler_config.toml")
	s.logger = logger.New(s.calendarConfig.Logger.Level, s.calendarConfig.Logger.Output)
	s.storage = storage.NewStorage(s.calendarConfig.Storage)
}

func (s *CalendarIntegrationSuite) SetupTest() {
	s.storage.Connect(context.Background())
	s.cleanupDB()
}

func (s *CalendarIntegrationSuite) TearDownTest() {
	s.cleanupDB()
	s.storage.Close(context.Background())
}

func (s *CalendarIntegrationSuite) cleanupDB() {
	s.storage.Truncate(context.Background())
}

func (s *CalendarIntegrationSuite) TestCrud() {
	addr := fmt.Sprintf("%v:%v", CalendarHost, s.calendarConfig.Endpoint.GRPCPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Suite.Require().NoError(err)

	client := pb.NewEventsClient(conn)
	createReq := &pb.NewEventRequest{
		Event: &pb.TransientEvent{
			Title:       "Event 1",
			Description: "some description",
			StartTime:   time.Now().Add(-24 * time.Hour).Unix(),
			EndTime:     time.Now().Add(-24 * time.Hour).Add(60 * time.Second).Unix(),
			OwnerEmail:  "user1@example.com",
			Notify:      "",
		},
	}

	scalarResponse, err := client.CreateEvent(context.Background(), createReq)
	s.Suite.Require().NoError(err, "create event failed")
	s.Suite.Require().NotNil(scalarResponse, "create response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "create response payload should not be nil")

	event := scalarResponse.GetEvent()

	getRequest := &pb.EventIdRequest{Id: event.Id}
	scalarResponse, err = client.GetEvent(context.Background(), getRequest)
	s.Suite.Require().NoError(err, "get event failed")
	s.Suite.Require().NotNil(scalarResponse, "get response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "get response payload should not be nil")

	updateReq := &pb.UpdateEventRequest{
		Id: event.Id,
		Event: &pb.TransientEvent{
			Title:       event.Title,
			Description: "new description",
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			OwnerEmail:  event.OwnerEmail,
			Notify:      event.Notify,
		},
	}

	scalarResponse, err = client.UpdateEvent(context.Background(), updateReq)
	s.Suite.Require().NoError(err, "update event failed")
	s.Suite.Require().NotNil(scalarResponse, "update response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "update response payload should not be nil")

	dateReq := &pb.DateRequest{
		Owner: "user1@example.com",
		Date:  time.Now().AddDate(0, 0, -7).Unix(),
	}
	vectorResponse, err := client.GetEventsForWeek(context.Background(), dateReq)
	s.Suite.Require().NoError(err, "get evetns for week failed")
	s.Suite.Require().NotNil(scalarResponse, "get evetns for week response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "get evetns for week response payload should not be nil")
	events := vectorResponse.GetEvents()
	s.Suite.Require().Len(events, 1, "wrong number of events")

	deleteRequest := &pb.EventIdRequest{Id: event.Id}
	_, err = client.DeleteEvent(context.Background(), deleteRequest)
	s.Suite.Require().NoError(err, "delete event failed")

	_, err = client.GetEvent(context.Background(), getRequest)
	s.Suite.Require().Error(err)
}

func (s *CalendarIntegrationSuite) TestSchedulerOldEventCleanup() {
	addr := fmt.Sprintf("%v:%v", CalendarHost, s.calendarConfig.Endpoint.GRPCPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Suite.Require().NoError(err)

	retentionPeriod, err := time.ParseDuration(s.scheduleConfig.Schedule.RetentionPeriod)
	s.Suite.Require().NoError(err, "wrong retention period")

	scanPeriod, err := time.ParseDuration(s.scheduleConfig.Schedule.Interval)
	s.Suite.Require().NoError(err, "wrong scan interval period")

	client := pb.NewEventsClient(conn)
	createReq := &pb.NewEventRequest{
		Event: &pb.TransientEvent{
			Title:       "Event Old 1",
			Description: "some description",
			StartTime:   time.Now().Add(-2 * retentionPeriod).Unix(),
			EndTime:     time.Now().Add(-2 * retentionPeriod).Add(30 * time.Second).Unix(),
			OwnerEmail:  "user1@example.com",
			Notify:      "",
		},
	}
	scalarResponse, err := client.CreateEvent(context.Background(), createReq)
	s.Suite.Require().NoError(err, "create event failed")
	s.Suite.Require().NotNil(scalarResponse, "create response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "create response payload should not be nil")
	event := scalarResponse.GetEvent()

	time.Sleep(scanPeriod * 2)

	getRequest := &pb.EventIdRequest{Id: event.Id}
	_, err = client.GetEvent(context.Background(), getRequest)
	s.Suite.Require().Error(err, "event should not exist")
}

func (s *CalendarIntegrationSuite) TestSenderSendNotification() {
	const shift = 20
	const notifyBefore = "10s"
	addr := fmt.Sprintf("%v:%v", CalendarHost, s.calendarConfig.Endpoint.GRPCPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Suite.Require().NoError(err)
	s.Suite.Require().NoError(err, "wrong retention period")

	client := pb.NewEventsClient(conn)
	createReq := &pb.NewEventRequest{
		Event: &pb.TransientEvent{
			Title:       "Event Notify 1",
			Description: "some description",
			StartTime:   time.Now().Add(shift * time.Second).Unix(),
			EndTime:     time.Now().Add(shift * time.Second).Unix(),
			OwnerEmail:  "user1@example.com",
			Notify:      notifyBefore,
		},
	}
	scalarResponse, err := client.CreateEvent(context.Background(), createReq)
	s.Suite.Require().NoError(err, "create event failed")
	s.Suite.Require().NotNil(scalarResponse, "create response should not be nil")
	s.Suite.Require().NotNil(scalarResponse.GetEvent(), "create response payload should not be nil")
	event := scalarResponse.GetEvent()

	time.Sleep(shift * time.Second)

	dbEvent, err := s.storage.GetEvent(context.Background(), event.Id)
	s.Suite.Require().NoError(err, "event not found")
	s.Suite.Require().True(dbEvent.Notified, "should be notified")
}
