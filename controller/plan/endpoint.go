package plan

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &npEndpoint{}
	plugin.RegisterEndpoint(e)
}

type npEndpoint struct {
	m *model.Model
}

func (a *npEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/plan", a.Create).Methods("POST")
	r.HandleFunc("/plan/{person}", a.Show).Methods("GET")
	r.HandleFunc("/plan/{person}", a.Delete).Methods("DELETE")
	r.HandleFunc("/plan/{person}", a.Update).Methods("PUT")
}

func (a *npEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for notification plan
}

func (a *npEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for notification plan
}

func (a *npEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
}

func (a *npEndpoint) Name() string {
	return "NotificationPlanEndpoint"
}
