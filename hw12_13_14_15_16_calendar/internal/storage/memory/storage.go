package memorystorage

import (
	"sync"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
)

type Storage struct {
	events map[string]storagecommon.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storagecommon.Event),
	}
}

func (s *Storage) Create(event storagecommon.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return storagecommon.ErrAlreadyExists
	}

	for _, e := range s.events {
		if e.UserID == event.UserID && isOverlapping(e, event) {
			return storagecommon.ErrConflictOverlap
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) GetByID(id string) (storagecommon.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id]
	if !ok {
		return storagecommon.Event{}, storagecommon.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) Update(event storagecommon.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.events[event.ID]
	if !exist {
		return storagecommon.ErrEventNotFound
	}

	for id, e := range s.events {
		if id != event.ID && e.UserID == event.UserID && isOverlapping(e, event) {
			return storagecommon.ErrConflictOverlap
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storagecommon.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) List() ([]storagecommon.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storagecommon.Event, 0, len(s.events))
	for _, v := range s.events {
		result = append(result, v)
	}
	return result, nil
}

func (s *Storage) ListByUser(userID string) ([]storagecommon.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storagecommon.Event, 0)
	for _, event := range s.events {
		if event.UserID == userID {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storagecommon.Event, 0)
	for _, event := range s.events {
		if event.UserID == userID && !event.EndTime.Before(from) && !event.StartTime.After(to) {
			result = append(result, event)
		}
	}
	return result, nil
}

func isOverlapping(a, b storagecommon.Event) bool {
	return a.StartTime.Before(b.EndTime) && b.StartTime.Before(a.EndTime)
}
