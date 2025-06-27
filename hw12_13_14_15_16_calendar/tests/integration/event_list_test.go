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

func TestListEvents_Day(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "user_1d"
	now := time.Now().UTC()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	createTestEvent(ctx, t, userID, "Today event", now.Unix(), now.Add(time.Hour).Unix())

	url := fmt.Sprintf("%s/events/range?userId=%s&from=%d&to=%d",
		calendarBaseURL, userID, startOfDay.Unix(), endOfDay.Unix())

	doRangeRequestAndVerify(ctx, t, url, func(events []Event) {
		assert.NotEmpty(t, events)
		for _, e := range events {
			assert.Equal(t, userID, e.UserID)
			assert.True(t, e.StartTime >= startOfDay.Unix() && e.StartTime <= endOfDay.Unix())
		}
	})
}

func TestListEvents_Week(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "user_1w"
	now := time.Now().UTC()

	var weekStart time.Time
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6)
	} else {
		weekStart = now.AddDate(0, 0, -int(now.Weekday())+1)
	}
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	weekEnd := weekStart.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)

	eventTime := weekStart.Add(12 * time.Hour)
	createTestEvent(ctx, t, userID, "Weekly event", eventTime.Unix(), eventTime.Add(time.Hour).Unix())

	url := fmt.Sprintf("%s/events/range?userId=%s&from=%d&to=%d",
		calendarBaseURL, userID, weekStart.Unix(), weekEnd.Unix())

	doRangeRequestAndVerify(ctx, t, url, func(events []Event) {
		assert.NotEmpty(t, events)
		for _, e := range events {
			assert.Equal(t, userID, e.UserID)
			assert.True(t, e.StartTime >= weekStart.Unix() && e.StartTime <= weekEnd.Unix())
		}
	})
}

func TestListEvents_Month(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "user_1m"
	now := time.Now().UTC()

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	eventTime := monthStart.Add(24 * time.Hour)
	createTestEvent(ctx, t, userID, "Monthly event", eventTime.Unix(), eventTime.Add(time.Hour).Unix())

	url := fmt.Sprintf("%s/events/range?userId=%s&from=%d&to=%d",
		calendarBaseURL, userID, monthStart.Unix(), monthEnd.Unix())

	doRangeRequestAndVerify(ctx, t, url, func(events []Event) {
		assert.NotEmpty(t, events)
		for _, e := range events {
			assert.Equal(t, userID, e.UserID)
			assert.True(t, e.StartTime >= monthStart.Unix() && e.StartTime <= monthEnd.Unix())
		}
	})
}

type Event struct {
	UserID    string `json:"userId"`
	Title     string `json:"title"`
	StartTime int64  `json:"startTime"`
}

type Response struct {
	Events []Event `json:"events"`
}

func doRangeRequestAndVerify(ctx context.Context, t *testing.T, url string, verifyFunc func([]Event)) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	t.Logf("Found %d events in range: %s", len(result.Events), url)

	verifyFunc(result.Events)
}
