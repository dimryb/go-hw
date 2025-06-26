package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSchedulerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	now := time.Now().UTC()
	startTime := now.Add(20 * time.Second).Unix()
	endTime := now.Add(1 * time.Hour).Unix()
	notifyBefore := int64(10)

	eventReq := CreateEventRequest{
		UserID:       "user_test_scheduler",
		Title:        "Integration Test Event",
		Description:  "Created by integration test for scheduler and sender",
		StartTime:    startTime,
		EndTime:      endTime,
		NotifyBefore: notifyBefore,
	}

	reqBody, _ := json.Marshal(eventReq)

	t.Run("CreateEvent", func(t *testing.T) {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", calendarBaseURL+"/event/create", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 Created")

		var createResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResp)
		assert.NoError(t, err)
		assert.Equal(t, "created", createResp["status"])
	})

	t.Run("WaitForNotificationInRabbitMQ", func(t *testing.T) {
		rmqClient, err := rmq.NewClient(rabbitURL, exchange)
		require.NoError(t, err, "Failed to create RMQ client")
		defer func() { _ = rmqClient.Close() }()

		msgs, err := rmqClient.Consume(exchange)
		require.NoError(t, err, "Failed to register a consumer")

		select {
		case msg := <-msgs:
			var notification Notification
			err := json.Unmarshal(msg, &notification)
			require.NoError(t, err)

			assert.Equal(t, "Integration Test Event", notification.Title)
			assert.Equal(t, "user_test_scheduler", notification.UserID)
			t.Logf("Received notification: %+v", notification)

		case <-time.After(60 * time.Second):
			t.Fatal("Timeout waiting for message in RabbitMQ queue")
		}
	})
}
