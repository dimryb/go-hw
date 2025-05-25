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
