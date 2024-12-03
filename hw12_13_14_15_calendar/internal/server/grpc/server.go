package internalgrpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/contracts"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application interface {
	CreateEvent(ctx context.Context, dto contracts.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, dto contracts.Event) (storage.Event, error)
	GetEvent(ctx context.Context, eventID string) (storage.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error
	ListEventsForDate(ctx context.Context, ownerEmail string, date int64) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, ownerEmail string, date int64) ([]storage.Event, error)
}

type Server struct {
	app Application
	pb.UnimplementedEventsServer
}

func (s *Server) toPersistedEvent(ev storage.Event) *pb.PersistedEvent {
	return &pb.PersistedEvent{
		Id:          ev.ID,
		Title:       ev.Title,
		StartTime:   ev.StartTime.Unix(),
		EndTime:     ev.EndTime.Unix(),
		Description: ev.Description,
		OwnerEmail:  ev.OwnerEmail,
		Notify:      ev.NotifyBefore,
	}
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.NewEventRequest) (*pb.ScalarEventResponse, error) {
	payload := req.GetEvent()
	if payload == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}

	event, err := s.app.CreateEvent(ctx, contracts.Event{
		Title:        payload.Title,
		StartTime:    payload.StartTime,
		EndTime:      payload.EndTime,
		Description:  payload.Description,
		OwnerEmail:   payload.OwnerEmail,
		NotifyBefore: payload.Notify,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ScalarEventResponse{
		Event: s.toPersistedEvent(event),
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.ScalarEventResponse, error) {
	payload := req.GetEvent()
	if payload == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}

	event, err := s.app.UpdateEvent(ctx, contracts.Event{
		ID:           payload.Id,
		Title:        payload.Title,
		StartTime:    payload.StartTime,
		EndTime:      payload.EndTime,
		Description:  payload.Description,
		OwnerEmail:   payload.OwnerEmail,
		NotifyBefore: payload.Notify,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ScalarEventResponse{
		Event: s.toPersistedEvent(event),
	}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.EventIdRequest) (*pb.ScalarEventResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	event, err := s.app.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.ScalarEventResponse{
		Event: s.toPersistedEvent(event),
	}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.EventIdRequest) (*empty.Empty, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	err := s.app.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) GetEventsForDay(ctx context.Context, req *pb.DateRequest) (*pb.VectorEventResponse, error) {
	owner := req.GetOwner()
	date := req.GetDate()
	if owner == "" {
		return nil, status.Error(codes.InvalidArgument, "owner is not specified")
	}

	result, err := s.app.ListEventsForDate(ctx, owner, date)
	if err != nil {
		return nil, err
	}

	persitedEvents := make([]*pb.PersistedEvent, 0)
	for _, e := range result {
		persitedEvents = append(persitedEvents, s.toPersistedEvent(e))
	}

	return &pb.VectorEventResponse{
		Events: persitedEvents,
	}, nil
}

func (s *Server) GetEventsForWeek(ctx context.Context, req *pb.DateRequest) (*pb.VectorEventResponse, error) {
	owner := req.GetOwner()
	date := req.GetDate()
	if owner == "" {
		return nil, status.Error(codes.InvalidArgument, "owner is not specified")
	}

	result, err := s.app.ListEventsForWeek(ctx, owner, date)
	if err != nil {
		return nil, err
	}

	persistedEvents := make([]*pb.PersistedEvent, 0)
	for _, e := range result {
		persistedEvents = append(persistedEvents, s.toPersistedEvent(e))
	}

	return &pb.VectorEventResponse{
		Events: persistedEvents,
	}, nil
}
