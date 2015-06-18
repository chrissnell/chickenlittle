package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

// NotifyPersonRequest is the JSON request sent by a client to trigger
// the notification of a single person.
type NotifyPersonRequest struct {
	Content string `json:"content"`
	UUID    string `json:"-"`
}

// NotifyPersonResponse is the response sent to the client in response
// to the NotifyPersonRequest.
type NotifyPersonResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

// NotifyPerson notifies a Person by looking up their NotificationPlan and sending it to the person notification engine.
func (a *Controller) NotifyPerson(w http.ResponseWriter, r *http.Request) {
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

	// TODO create an EscalationPlan object here, see notification.NotificationForPerson()
	//req.Plan, err = a.m.GetNotificationPlan(username)
	//if err != nil {
	//	// res.Error = err.Error()
	//	// errjson, _ := json.Marshal(res)
	//	http.Error(w, err.Error(), http.StatusNotFound)
	//	return
	//}

	// Assign a UUID to this notification.  The UUID is used to track notifications-in-progress (NIP) and to stop
	// them when requested.
	uuid.SwitchFormat(uuid.CleanHyphen)
	req.UUID = uuid.NewV4().String()

	//// Send our NotificationRequest to the notification engine
	//a.n.EnqueueNotification(&req)

	res = NotifyPersonResponse{
		Message:  "Notification initiated",
		Content:  req.Content,
		UUID:     req.UUID,
		Username: username,
	}

	json.NewEncoder(w).Encode(res)

}
