package http

import (
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

func TestGetEventByID(t *testing.T) {
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
	_, err = testApp.Storage.Create(initialEvent)
	require.NoError(t, err)

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/event/get?id=event123", nil)
	w := httptest.NewRecorder()

	testApp.Server.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response internalhttp.EventResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, initialEvent.ID, response.ID)
	assert.Equal(t, initialEvent.UserID, response.UserID)
	assert.Equal(t, initialEvent.Title, response.Title)
	assert.Equal(t, initialEvent.Description, response.Description)
	assert.Equal(t, initialEvent.StartTime.Unix(), response.StartTime)
	assert.Equal(t, initialEvent.EndTime.Unix(), response.EndTime)
	assert.Equal(t, int64(initialEvent.NotifyBefore), response.NotifyBefore)
}
