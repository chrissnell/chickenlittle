package notification

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// StopNotification stops a notification-in-progress (NIP) by sending the UUID to the notification engine
func (a *notificationEndpoint) StopNotification(w http.ResponseWriter, r *http.Request) {
	var res Response

	vars := mux.Vars(r)
	id := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !a.n.IsNotification(id) {
		res = Response{
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

	res = Response{
		Message: "Attempting to terminate notification",
		UUID:    id,
	}

	json.NewEncoder(w).Encode(res)
}

// StopNotificationClick provides a simple GET-able endpoint to stop notifications when a link is clicked in an email client
func (a *notificationEndpoint) StopNotificationClick(w http.ResponseWriter, r *http.Request) {

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
