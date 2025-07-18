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

func TestListEvents(t *testing.T) {
	testApp := tests.NewTestAppForCalendar()
	err := testApp.Setup()
	require.NoError(t, err)
	defer testApp.Teardown()

	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	eventsToCreate := []storagecommon.Event{
		{
			ID:           "event1",
			UserID:       "user123",
			Title:        "Event 1",
			Description:  "Desc 1",
			StartTime:    now,
			EndTime:      now.Add(time.Hour),
			NotifyBefore: 600,
		},
		{
			ID:           "event2",
			UserID:       "user456",
			Title:        "Event 2",
			Description:  "Desc 2",
			StartTime:    now.Add(2 * time.Hour),
			EndTime:      now.Add(3 * time.Hour),
			NotifyBefore: 900,
		},
	}

	for _, e := range eventsToCreate {
		_, err := testApp.Storage.Create(e)
		require.NoError(t, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/events/list", nil)
	w := httptest.NewRecorder()

	testApp.Server.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response internalhttp.ListEventsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	actualEvents := make([]storagecommon.Event, 0, len(response.Events))
	for _, item := range response.Events {
		actualEvents = append(actualEvents, storagecommon.Event{
			ID:           item.ID,
			UserID:       item.UserID,
			Title:        item.Title,
			Description:  item.Description,
			StartTime:    time.Unix(item.StartTime, 0).UTC(),
			EndTime:      time.Unix(item.EndTime, 0).UTC(),
			NotifyBefore: int(item.NotifyBefore),
		})
	}

	assert.ElementsMatch(t, eventsToCreate, actualEvents)
}
