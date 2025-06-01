package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

type CalendarHandlers struct {
	app    Application
	logger Logger
}

func NewCalendarHandlers(app Application, logger Logger) *CalendarHandlers {
	return &CalendarHandlers{app, logger}
}

func (h *CalendarHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.StartTime >= req.EndTime {
		http.Error(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	event := types.Event{
		UserID:       req.UserID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		NotifyBefore: int(req.NotifyBefore),
	}

	ctx := r.Context()

	if err := h.app.CreateEvent(ctx, event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "created"}); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}
