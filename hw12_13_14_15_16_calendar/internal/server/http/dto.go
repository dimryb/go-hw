package internalhttp

type CreateEventRequest struct {
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	NotifyBefore int64  `json:"notifyBefore"`
}

type UpdateEventRequest struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	NotifyBefore int64  `json:"notifyBefore"`
}

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
