package people

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NotifyPerson notifies a Person by looking up their NotificationPlan and sending it to the person notification engine.
func (a *personEndpoint) Notify(w http.ResponseWriter, r *http.Request) {
	var res NotificationResponse
	var req NotificationRequest

	vars := mux.Vars(r)
	username := vars["person"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*20))
	// If something went wrong, return an error in the JSON response
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = r.Body.Close()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Error = "Error unmarshaling request: " + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	n, err := a.m.GetNotificationForPerson(username, req.Summary, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send our NotificationRequest to the notification engine
	a.n.EnqueueNotification(&n)
	log.Printf("Enqueued notification: %v", n)

	res = NotificationResponse{
		Message:  "Notification initiated",
		Content:  n.Message(),
		UUID:     n.ID(),
		Username: username,
	}

	json.NewEncoder(w).Encode(res)
}
