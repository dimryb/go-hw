package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEventsByUserInRange(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "user_1"
	now := time.Now().UTC()
	from := now.Unix()
	to := now.Add(24 * time.Hour).Unix()

	createTestEvent(ctx, t, userID, "Event 1", from+3600, from+7200)
	createTestEvent(ctx, t, userID, "Event 2", from+10800, from+14400)

	url := fmt.Sprintf("%s/events/range?userId=%s&from=%d&to=%d", calendarBaseURL, userID, from, to)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type Event struct {
		UserID    string `json:"userId"`
		Title     string `json:"title"`
		StartTime int64  `json:"startTime"`
	}

	type Response struct {
		Events []Event `json:"events"`
	}

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	t.Logf("Found %d events for user %s in range", len(result.Events), userID)

	assert.NotEmpty(t, result.Events, "Expected at least one event")

	for _, e := range result.Events {
		assert.Equal(t, userID, e.UserID)
		assert.True(t, e.StartTime >= from && e.StartTime <= to)
	}
}

func createTestEvent(ctx context.Context, t *testing.T, userID, title string, start, end int64) {
	t.Helper()
	req := CreateEventRequest{
		UserID:    userID,
		Title:     title,
		StartTime: start,
		EndTime:   end,
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", calendarBaseURL+"/event/create", bytes.NewBuffer(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
}
