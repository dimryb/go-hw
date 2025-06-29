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

func TestNotificationIsSent(t *testing.T) {
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

	t.Run("WaitForNotificationStatus", func(t *testing.T) {
		rmqClient, err := rmq.NewClient(rabbitURL, "notification_status")
		require.NoError(t, err, "Failed to create RMQ client for status")
		defer func() { _ = rmqClient.Close() }()

		msgs, err := rmqClient.Consume("notification_status")
		require.NoError(t, err, "Failed to register a consumer for status")

		select {
		case msg := <-msgs:
			var status rmq.NotificationStatus
			err := json.Unmarshal(msg, &status)
			require.NoError(t, err)

			assert.Equal(t, "user_test_scheduler", status.UserID)
			assert.Equal(t, "delivered", status.Status)
			assert.NotZero(t, status.Timestamp)
			t.Logf("Received notification status: %+v", status)

		case <-time.After(60 * time.Second):
			t.Fatal("Timeout waiting for status message in RabbitMQ queue")
		}
	})
}
