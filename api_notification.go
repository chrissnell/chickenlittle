package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/chrissnell/chickenlittle/ne"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

type NotifyPersonRequest struct {
	Content string            `json:"content"`
	Plan    *NotificationPlan `json:"-"`
	step    int               `json:"-"`
}

func (n *NotifyPersonRequest) NextStep() (ne.NotificationStep, bool) {
	last := false
	if n.step < len(n.Plan.Steps)-1 {
		n.step++
	} else {
		last = true
	}
	return n.Plan.Steps[n.step], last
}

func (n *NotifyPersonRequest) ID() string {
	return n.Plan.ID.String()
}

func (n *NotifyPersonRequest) Message() string {
	return n.Content
}

func (n *NotifyPersonRequest) Subject() string {
	return "Chicken Little notification"
}

type NotifyTeamRequest struct {
	Content string `json:"content"`
	Team    *Team  `json:"-"`
	Id      uuid.UUID
	step    int `json:"-"`
}

func (n *NotifyTeamRequest) NextStep() (ne.NotificationStep, bool) {
	// TODO implement escalation logic ...
	return NotificationStep{}, true
}

func (n *NotifyTeamRequest) ID() string {
	return n.Id.String()
}

func (n *NotifyTeamRequest) Message() string {
	return n.Content
}

func (n *NotifyTeamRequest) Subject() string {
	return "Chicken Little notification"
}

type NotifyPersonResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

type NotifyTeamResponse struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// NotifyTeam notifies a team by looking up the escalation order and sending it to the team notification engine.
func NotifyTeam(w http.ResponseWriter, r *http.Request) {
	var res NotifyTeamResponse
	var req NotifyTeamRequest

	vars := mux.Vars(r)
	name := vars["name"]

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

	req.Team, err = c.GetTeam(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID to this notification. The UUID is used to track notifications-in-progress and to stop
	// them when requested.
	uuid.SwitchFormat(uuid.CleanHyphen)
	req.Id = uuid.NewV4()

	// Send our notification request to the notification engine
	c.Notify.EnqueueNotification(&req)

	res = NotifyTeamResponse{
		Message: "Notification initated",
		Content: req.Content,
		UUID:    req.ID(),
		Name:    name,
	}

	json.NewEncoder(w).Encode(res)
}

// NotifyPerson notifies a Person by looking up their NotificationPlan and sending it to the person notification engine.
func NotifyPerson(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse
	var req NotifyPersonRequest

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

	req.Plan, err = c.GetNotificationPlan(username)
	if err != nil {
		// res.Error = err.Error()
		// errjson, _ := json.Marshal(res)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID to this notification.  The UUID is used to track notifications-in-progress (NIP) and to stop
	// them when requested.
	uuid.SwitchFormat(uuid.CleanHyphen)
	req.Plan.ID = uuid.NewV4()

	// Send our NotificationRequest to the notification engine
	c.Notify.EnqueueNotification(&req)

	res = NotifyPersonResponse{
		Message:  "Notification initiated",
		Content:  req.Content,
		UUID:     req.Plan.ID.String(),
		Username: req.Plan.Username,
	}

	json.NewEncoder(w).Encode(res)

}

// Stop a notification-in-progress (NIP) by sending the UUID to the notification engine
func StopNotification(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	id := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !c.Notify.IsNotification(id) {
		res = NotifyPersonResponse{
			Error: "No active notifications for this UUID",
			UUID:  id,
		}
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	// Attempt to stop the notification by sending the UUID to the notification engine
	c.Notify.CancelNotification(id)

	// TO DO: make sure that this is a valid UUID and obtain
	//        confirmation of deletion

	res = NotifyPersonResponse{
		Message: "Attempting to terminate notification",
		UUID:    id,
	}

	json.NewEncoder(w).Encode(res)
}

// A simple GET-able endpoint to stop notifications when a link is clicked in an email client
func StopNotificationClick(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["uuid"]

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if !c.Notify.IsNotification(id) {
		http.Error(w, fmt.Sprintf("UUID %s not found.", id), http.StatusNotFound)
		return
	}

	// Attempt to stop the notification by sending the UUID to the notification engine
	c.Notify.CancelNotification(id)

	fmt.Fprintln(w, "<html><body><b>Thank you!</b><br><br>Chicken Little has received your acknowledgement and you will no longer be notified with this message.</body></html>")
}
