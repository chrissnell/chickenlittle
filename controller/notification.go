package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NotificationRequest is the JSON request sent by a client to trigger
// to escalating notification of a whole team. Implements the Notification interface.
type NotificationRequest struct {
	Summary string `json:"summary"` // a summary or subject of the notification. Currenlty not used in all integrations. Optional.
	Content string `json:"content"` // the notification content. mandatory.
}

// StopNotification stops a notification-in-progress (NIP) by sending the UUID to the notification engine
func (a *Controller) StopNotification(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	id := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !a.n.IsNotification(id) {
		res = NotifyPersonResponse{
			Error: "No active notifications for this UUID",
			UUID:  id,
		}
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	// Attempt to stop the notification by sending the UUID to the notification engine
	a.n.CancelNotification(id)

	// TO DO: make sure that this is a valid UUID and obtain
	//        confirmation of deletion

	res = NotifyPersonResponse{
		Message: "Attempting to terminate notification",
		UUID:    id,
	}

	json.NewEncoder(w).Encode(res)
}

// StopNotificationClick provides a simple GET-able endpoint to stop notifications when a link is clicked in an email client
func (a *Controller) StopNotificationClick(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["uuid"]

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if !a.n.IsNotification(id) {
		http.Error(w, fmt.Sprintf("UUID %s not found.", id), http.StatusNotFound)
		return
	}

	// Attempt to stop the notification by sending the UUID to the notification engine
	a.n.CancelNotification(id)

	fmt.Fprintln(w, "<html><body><b>Thank you!</b><br><br>Chicken Little has received your acknowledgement and you will no longer be notified with this message.</body></html>")
}
