package notification

// Request is the JSON request sent by a client to trigger
// to escalating notification of a whole team. Implements the Notification interface.
type Request struct {
	Summary string `json:"summary"` // a summary or subject of the notification. Currenlty not used in all integrations. Optional.
	Content string `json:"content"` // the notification content. mandatory.
}

// Response is the response sent to the client in response
// to a NotificationRequest
type Response struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}
