package team

import "github.com/chrissnell/chickenlittle/model"

// TeamsResponse is a struct to return one or more teams as JSON
type TeamsResponse struct {
	Teams   []model.Team `json:"teams"`
	Message string       `json:"message"`
	Error   string       `json:"error"`
}

// NotificationRequest is the JSON request sent by a client to trigger
// to escalating notification of a whole team. Implements the Notification interface.
type NotificationRequest struct {
	Summary string `json:"summary"` // a summary or subject of the notification. Currenlty not used in all integrations. Optional.
	Content string `json:"content"` // the notification content. mandatory.
}

// NotifyTeamResponse is the response sent to the client in response
// to the NotifyTeamRequest.
type NotifyTeamResponse struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
