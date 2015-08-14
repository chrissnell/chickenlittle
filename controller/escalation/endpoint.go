package escalation

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &escalationEndpoint{}
	plugin.RegisterEndpoint(e)
}

type escalationEndpoint struct {
	m *model.Model
}

func (a *escalationEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/escalation", a.Create).Methods("POST")
	r.HandleFunc("/escalation/{plan}", a.Show).Methods("GET")
	r.HandleFunc("/escalation/{plan}", a.Delete).Methods("DELETE")
	r.HandleFunc("/escalation/{plan}", a.Update).Methods("PUT")
}

func (a *escalationEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for escalation plans
}

func (a *escalationEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for escalation plans
}

func (a *escalationEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
}

func (a *escalationEndpoint) Name() string {
	return "EscalationEndpoint"
}
