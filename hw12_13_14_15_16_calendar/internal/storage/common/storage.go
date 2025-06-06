package storagecommon

import "time"

//go:generate mockgen -source=storage.go -package=mocks -destination=../../../mocks/mock_storage.go
type EventStorage interface {
	Create(event Event) (string, error)
	Update(event Event) error
	Delete(id string) error

	GetByID(id string) (Event, error)
	List() ([]Event, error)
	ListByUser(userID string) ([]Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]Event, error)
}
