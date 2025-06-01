package types

import "time"

type Event struct {
	ID           string
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	UserID       string
	NotifyBefore int
}
