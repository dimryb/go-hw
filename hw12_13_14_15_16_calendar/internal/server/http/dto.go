package internalhttp

// CreateEventRequest represents the request to create an event.
// @Description Represents the request to create an event.
type CreateEventRequest struct {
	UserID       string `json:"userId" example:"id1234"`
	Title        string `json:"title" example:"Team Meeting"`
	Description  string `json:"description" example:"Discuss project roadmap"`
	StartTime    int64  `json:"startTime" example:"1717290000"`
	EndTime      int64  `json:"endTime" example:"1717293600"`
	NotifyBefore int64  `json:"notifyBefore" example:"600"`
}

// UpdateEventRequest represents the request to update an existing event.
// @Description Represents the request to update an existing event.
type UpdateEventRequest struct {
	ID           string `json:"id" example:"event5678"`
	UserID       string `json:"userId" example:"id1234"`
	Title        string `json:"title" example:"Team Meeting Updated"`
	Description  string `json:"description" example:"Updated description"`
	StartTime    int64  `json:"startTime" example:"1717290000"`
	EndTime      int64  `json:"endTime" example:"1717293600"`
	NotifyBefore int64  `json:"notifyBefore" example:"700"`
}

// EventResponse represents an event returned by the API.
// @Description Represents an event returned by the API.
type EventResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	NotifyBefore int64  `json:"notifyBefore"`
}

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

// CreateEventResponse represents a successful event creation response.
// @Description Represents a successful event creation response.
type CreateEventResponse struct {
	Status string             `json:"status"`
	Event  CreateEventRequest `json:"event"`
}
