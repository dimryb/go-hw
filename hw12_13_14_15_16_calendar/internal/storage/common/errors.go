package storagecommon

import "fmt"

var (
	ErrEventNotFound   = fmt.Errorf("event not found")
	ErrDateBusy        = fmt.Errorf("the selected time is already busy")
	ErrInvalidEvent    = fmt.Errorf("invalid event data")
	ErrAlreadyExists   = fmt.Errorf("event already exists")
	ErrConflictOverlap = fmt.Errorf("event overlaps with another event")
)
