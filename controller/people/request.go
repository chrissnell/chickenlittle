package people

import "github.com/chrissnell/chickenlittle/model"

// Response is the JSON response
type Response struct {
	People  []model.Person `json:"people"`
	Message string         `json:"message"`
	Error   string         `json:"error"`
}

// NotificationRequest is the JSON request sent by a client to trigger
// to escalating notification of a whole team. Implements the Notification interface.
type NotificationRequest struct {
	Summary string `json:"summary"` // a summary or subject of the notification. Currenlty not used in all integrations. Optional.
	Content string `json:"content"` // the notification content. mandatory.
}

// NotificationResponse is the response sent to the client in response
// to the NotifyPersonRequest.
type NotificationResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}
