package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CreateEventRequest struct {
	ID           string `json:"id,omitempty"`
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	NotifyBefore int64  `json:"notifyBefore"`
}

type CreateEventResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func TestCreateEvent_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now().UTC()
	startTime := now.Add(1 * time.Hour).Unix()
	endTime := startTime + 3600

	req := CreateEventRequest{
		UserID:       "user_1",
		Title:        "Meeting",
		Description:  "Team meeting",
		StartTime:    startTime,
		EndTime:      endTime,
		NotifyBefore: 600,
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", calendarBaseURL+"/event/create", bytes.NewBuffer(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result CreateEventResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "created", result.Status)
	assert.NotEmpty(t, result.ID)
}
