package team

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// NotifyTeam notifies a team by looking up the escalation order and sending it to the team notification engine.
func (a *teamEndpoint) Notify(w http.ResponseWriter, r *http.Request) {
	var res NotifyTeamResponse
	var req NotificationRequest

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
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	n, err := a.m.GetNotificationForTeam(name, req.Summary, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send our notification request to the notification engine
	a.n.EnqueueNotification(&n)

	res = NotifyTeamResponse{
		Message: "Notification initated",
		Content: n.Message(),
		UUID:    n.ID(),
		Name:    name,
	}

	json.NewEncoder(w).Encode(res)
	return
}
