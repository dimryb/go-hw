package storagecommon

import "time"

type Event struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	StartTime    time.Time `db:"start_time"`
	EndTime      time.Time `db:"end_time"`
	Description  string    `db:"description"`
	UserID       string    `db:"user_id"`
	NotifyBefore int       `db:"notify_before"`
}

func (e Event) With(fn func(Event) Event) Event {
	return fn(e)
}

func (e Event) WithID(id string) Event {
	e.ID = id
	return e
}
