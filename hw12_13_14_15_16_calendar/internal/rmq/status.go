package rmq

import "time"

type NotificationStatus struct {
	NotificationID string    `json:"notificationId"`
	EventID        string    `json:"eventId"`
	UserID         string    `json:"userId"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
}
