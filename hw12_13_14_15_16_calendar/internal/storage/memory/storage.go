package memorystorage

import (
	"sync"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Create(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return storage.ErrAlreadyExists
	}

	for _, e := range s.events {
		if e.UserID == event.UserID && isOverlapping(e, event) {
			return storage.ErrConflictOverlap
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) GetByID(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id]
	if !ok {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) Update(event storage.Event) error {
	panic("not implemented")
}

func (s *Storage) Delete(id string) error {
	panic("not implemented")
}

func (s *Storage) List() ([]storage.Event, error) {
	panic("not implemented")
}

func (s *Storage) ListByUser(userID string) ([]storage.Event, error) {
	panic("not implemented")
}

func (s *Storage) ListByUserInRange(userID string, from, to time.Time) ([]storage.Event, error) {
	panic("not implemented")
}

func isOverlapping(a, b storage.Event) bool {
	return a.StartTime.Before(b.EndTime) && b.StartTime.Before(a.EndTime)
}
