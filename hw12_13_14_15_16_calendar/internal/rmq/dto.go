package rmq

type Notification struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
	Time        string `json:"time"`
	NotifyAt    string `json:"notifyAt"`
}
