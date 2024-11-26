package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) Connect(_ context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	// TODO
	return nil
}

func (s *Storage) AddEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[event.ID]
	if !ok {
		return customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", event.ID)}
	}
	s.events[event.ID] = event

	return nil
}

func (s *Storage) GetEvent(_ context.Context, eventID string) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[eventID]
	if !ok {
		return storage.Event{}, customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", eventID)}
	}
	return s.events[eventID], nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[eventID]
	if !ok {
		return customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", eventID)}
	}

	delete(s.events, eventID)
	return nil
}

func (s *Storage) ListEventsForPeriod(
	_ context.Context,
	ownerEmail string,
	startDate,
	endDate time.Time,
) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ev := range s.events {
		if ev.OwnerEmail != ownerEmail {
			continue
		}

		if (startDate.Before(ev.StartTime) || startDate.Equal(ev.StartTime)) &&
			(endDate.After(ev.StartTime) || endDate.Equal(ev.StartTime)) {
			result = append(result, ev)
		}
	}

	return result, nil
}
