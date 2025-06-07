package interfaces

import (
	"time"

	storagecommon "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
)

type Storage interface {
	Create(event storagecommon.Event) (string, error)
	Update(event storagecommon.Event) error
	Delete(id string) error

	GetByID(id string) (storagecommon.Event, error)
	List() ([]storagecommon.Event, error)
	ListByUser(userID string) ([]storagecommon.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error)
}
