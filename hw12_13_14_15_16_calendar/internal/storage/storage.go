package storage

import "time"

type EventStorage interface {
	Create(event Event) error
	Update(event Event) error
	Delete(id string) error

	GetByID(id string) (Event, error)
	List() ([]Event, error)
	ListByUser(userID string) ([]Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]Event, error)
}
