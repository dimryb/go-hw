package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	testApp := tests.NewTestAppForCalendar()
	err := testApp.Setup()
	require.NoError(t, err)
	defer testApp.Teardown()

	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	event := internalhttp.CreateEventRequest{
		UserID:       "user123",
		Title:        "Team Meeting",
		Description:  "Discuss roadmap",
		StartTime:    now.Unix(),
		EndTime:      now.Add(time.Hour).Unix(),
		NotifyBefore: 600,
	}

	body, err := json.Marshal(event)
	require.NoError(t, err)

	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/event/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	testApp.Server.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response internalhttp.CreateEventResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response.Status)
	// TODO: assert.NotEmpty(t, response.Id)

	list, err := testApp.Storage.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	if len(list) > 0 {
		created := list[0]
		assert.Equal(t, event.UserID, created.UserID)
		assert.Equal(t, event.Title, created.Title)
		assert.Equal(t, event.Description, created.Description)
		assert.Equal(t, time.Unix(event.StartTime, 0), created.StartTime)
		assert.Equal(t, time.Unix(event.EndTime, 0), created.EndTime)
		assert.Equal(t, int(event.NotifyBefore), created.NotifyBefore)
	}
}
