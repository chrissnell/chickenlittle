package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

// NotifyTeamRequest is the JSON request sent by a client to trigger
// to escalating notification of a whole team. Implements the Notification interface.
type NotifyTeamRequest struct {
	UUID    string                `json:"-"`
	Content string                `json:"content"`
	Plan    *model.EscalationPlan `json:"-"`
	Team    *model.Team           `json:"-"`
	member  int                   // current team member being notified, necessary for EscalationMethod = NotifyNextInRotation
	step    int                   // current step in the escalation plan
	ms      int                   // current notification step in the current members notification plan
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

// NextStep returns the next step of an escalation chain
func (n *NotifyTeamRequest) NextStep() (notification.NotificationStep, bool) {
	var step notification.NotificationStep
	var last bool

	// TODO not yet implemented
	//es := n.Plan.Steps[n.step]
	//if es.Method == model.NotifyOnDuty {
	//	// TODO ...
	//} else if es.Method == model.NotifyNextInRotation {
	//	n.member++
	//} else if es.Method == model.NotifyOtherPerson {
	//	// TODO ...
	//} else if es.Method == model.NotifyWebhook {
	//	step = notification.NotificationStep{
	//		Method:            es.Target, // TODO ensure Target starts with http or https
	//		NotifyEveryPeriod: es.TimeBeforeEscalation,
	//		NotifyUntilPeriod: es.TimeBeforeEscalation,
	//	}
	//} else if es.Method == model.NotifyEmail {
	//	step = notification.NotificationStep{
	//		Method:            "email://" + es.Target,
	//		NotifyEveryPeriod: es.TimeBeforeEscalation,
	//		NotifyUntilPeriod: es.TimeBeforeEscalation,
	//	}
	//}

	if n.step >= len(n.Plan.Steps)-1 {
		last = true
	}

	return step, last
}

// ID will return the UUID of this notification request
func (n *NotifyTeamRequest) ID() string {
	return n.UUID
}

// Message will return the message of this notification request
func (n *NotifyTeamRequest) Message() string {
	return n.Content
}

// Subject returns the default notification subject
func (n *NotifyTeamRequest) Subject() string {
	return "Chicken Little notification"
}

// Stopper will return the stop chan for this notification request
func (n *NotifyTeamRequest) Stopper() chan<- struct{} {
	// TODO implement
	return make(chan struct{})
}

// NotifyTeam notifies a team by looking up the escalation order and sending it to the team notification engine.
func (a *Controller) NotifyTeam(w http.ResponseWriter, r *http.Request) {
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

	req.Plan, err = a.m.GetEscalationPlan(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	req.Team, err = a.m.GetTeam(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID to this notification. The UUID is used to track notifications-in-progress and to stop
	// them when requested.
	uuid.SwitchFormat(uuid.CleanHyphen)
	req.UUID = uuid.NewV4().String()

	// Send our notification request to the notification engine
	a.n.EnqueueNotification(&req)

	res = NotifyTeamResponse{
		Message: "Notification initated",
		Content: req.Content,
		UUID:    req.UUID,
		Name:    name,
	}

	json.NewEncoder(w).Encode(res)
}
