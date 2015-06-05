package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

type NotifyPersonResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := c.GetNotificationPlan(username)
	if err != nil {
		// res.Error = err.Error()
		// errjson, _ := json.Marshal(res)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID
	uuid.SwitchFormat(uuid.CleanHyphen)
	p.ID = uuid.NewV4()

	planChan <- p

	res = NotifyPersonResponse{
		Message:  "Notification initiated",
		UUID:     p.ID.String(),
		Username: p.Username,
	}

	json.NewEncoder(w).Encode(res)

}

func StopNotification(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	id := vars["uuid"]

	// Attempt to stop the notification by sending the UUID to the notification engine
	stopChan <- id

	// TO DO: make sure that this is a valid UUID and obtain
	//        confirmation of deletion

	res = NotifyPersonResponse{
		Message: "Attempting to terminate notification",
		UUID:    id,
	}

	json.NewEncoder(w).Encode(res)
}
