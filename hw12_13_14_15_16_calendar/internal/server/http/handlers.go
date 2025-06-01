package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func toEventResponse(event types.Event) EventResponse {
	return EventResponse{
		ID:           event.ID,
		UserID:       event.UserID,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    event.StartTime.Unix(),
		EndTime:      event.EndTime.Unix(),
		NotifyBefore: int64(event.NotifyBefore),
	}
}

func (h *CalendarHandlers) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello, world!"))
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new calendar event
// @Tags events
// @Accept json
// @Produce json
// @Param event body CreateEventRequest true "Event data"
// @Success 201 {object} CreateEventResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /event/create [post].
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

func (h *CalendarHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	event := types.Event{
		ID:           req.ID,
		UserID:       req.UserID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		NotifyBefore: int(req.NotifyBefore),
	}

	ctx := r.Context()

	if err := h.app.UpdateEvent(ctx, event); err != nil {
		h.logger.Errorf("Failed to update event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "updated"}); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}

func (h *CalendarHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.app.DeleteEvent(ctx, id); err != nil {
		h.logger.Errorf("Failed to delete event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to delete event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "deleted"}); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}

func (h *CalendarHandlers) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	event, err := h.app.GetEventByID(ctx, id)
	if err != nil {
		h.logger.Errorf("Failed to get event by ID: %v", err)
		http.Error(w, fmt.Sprintf("Event not found: %v", err), http.StatusNotFound)
		return
	}

	resp := toEventResponse(event)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}

func (h *CalendarHandlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := h.app.ListEvents(ctx)
	if err != nil {
		h.logger.Errorf("Failed to list events: %v", err)
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	response := ListEventsResponse{}
	for _, e := range events {
		response.Events = append(response.Events, toEventResponse(e))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}

func (h *CalendarHandlers) ListEventsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	events, err := h.app.ListEventsByUser(ctx, userID)
	if err != nil {
		h.logger.Errorf("Failed to list events for user: %v", err)
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	response := ListEventsResponse{}
	for _, e := range events {
		response.Events = append(response.Events, toEventResponse(e))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}

func (h *CalendarHandlers) ListEventsByUserInRange(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	fromUnix, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid from timestamp", http.StatusBadRequest)
		return
	}
	toUnix, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid to timestamp", http.StatusBadRequest)
		return
	}

	from := time.Unix(fromUnix, 0)
	to := time.Unix(toUnix, 0)

	ctx := r.Context()
	events, err := h.app.ListEventsByUserInRange(ctx, userID, from, to)
	if err != nil {
		h.logger.Errorf("Failed to list events in range: %v", err)
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	response := ListEventsResponse{}
	for _, e := range events {
		response.Events = append(response.Events, toEventResponse(e))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
	}
}
