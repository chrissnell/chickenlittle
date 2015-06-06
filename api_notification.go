package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

type NotificationRequest struct {
	Content string            `json:"content"`
	Plan    *NotificationPlan `json:"-"`
}

type NotifyPersonResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse
	var req NotificationRequest

	vars := mux.Vars(r)
	username := vars["person"]

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

	req.Plan, err = c.GetNotificationPlan(username)
	if err != nil {
		// res.Error = err.Error()
		// errjson, _ := json.Marshal(res)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID
	uuid.SwitchFormat(uuid.CleanHyphen)
	req.Plan.ID = uuid.NewV4()

	planChan <- &req

	res = NotifyPersonResponse{
		Message:  "Notification initiated",
		Content:  req.Content,
		UUID:     req.Plan.ID.String(),
		Username: req.Plan.Username,
	}

	json.NewEncoder(w).Encode(res)

}

func StopNotification(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	id := vars["uuid"]

	if _, exists := NIP.Stoppers[id]; !exists {
		res = NotifyPersonResponse{
			Error: "No active notifications for this UUID",
			UUID:  id,
		}
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

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
