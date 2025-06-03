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
	storagecommon "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateEvent(t *testing.T) {
	testApp := tests.NewTestAppForCalendar()
	err := testApp.Setup()
	require.NoError(t, err)
	defer testApp.Teardown()

	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	initialEvent := storagecommon.Event{
		ID:           "event123",
		UserID:       "user123",
		Title:        "Old Title",
		Description:  "Old Description",
		StartTime:    now,
		EndTime:      now.Add(time.Hour),
		NotifyBefore: 600,
	}
	id, err := testApp.Storage.Create(initialEvent)
	require.NoError(t, err)

	updateReq := internalhttp.UpdateEventRequest{
		ID:           id,
		UserID:       "user123",
		Title:        "New Title",
		Description:  "New Description",
		StartTime:    now.Add(2 * time.Hour).Unix(),
		EndTime:      now.Add(3 * time.Hour).Unix(),
		NotifyBefore: 900,
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/event/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testApp.Server.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response internalhttp.UpdateEventResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "updated", response.Status)
	assert.Equal(t, "event123", response.ID)

	updatedEvent, err := testApp.Storage.GetByID("event123")
	require.NoError(t, err)

	assert.Equal(t, updateReq.Title, updatedEvent.Title)
	assert.Equal(t, updateReq.Description, updatedEvent.Description)
	assert.Equal(t, time.Unix(updateReq.StartTime, 0), updatedEvent.StartTime)
	assert.Equal(t, time.Unix(updateReq.EndTime, 0), updatedEvent.EndTime)
	assert.Equal(t, int(updateReq.NotifyBefore), updatedEvent.NotifyBefore)
}
