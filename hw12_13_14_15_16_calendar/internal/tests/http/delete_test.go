package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	storagecommon "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteEvent(t *testing.T) {
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

	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", "/event/delete?id=event123", nil)
	w := httptest.NewRecorder()

	testApp.Server.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = testApp.Storage.GetByID("event123")
	require.Error(t, err)
	assert.ErrorIs(t, err, storagecommon.ErrEventNotFound)
}
