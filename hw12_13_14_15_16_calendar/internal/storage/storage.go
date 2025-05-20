package storage

import (
	"fmt"
	"time"
)

type EventStorage interface {
	Create(event Event) error
	Update(event Event) error
	Delete(id string) error

	GetByID(id string) (Event, error)
	List() ([]Event, error)
	ListByUser(userID string) ([]Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]Event, error)
}

var (
	ErrEventNotFound   = fmt.Errorf("event not found")
	ErrDateBusy        = fmt.Errorf("the selected time is already busy")
	ErrInvalidEvent    = fmt.Errorf("invalid event data")
	ErrAlreadyExists   = fmt.Errorf("event already exists")
	ErrConflictOverlap = fmt.Errorf("event overlaps with another event")
)
