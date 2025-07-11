package integration

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type Notification struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
	Time        string `json:"time"`
	NotifyAt    string `json:"notifyAt"`
}

const (
	calendarBaseURL = "http://calendar-app:8080"
	rabbitURL       = "amqp://guest:guest@rabbitmq:5672/"
	exchange        = "notifications"
)

func TestCalendarHealth(t *testing.T) {
	url := "http://calendar-app:8080/"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}
