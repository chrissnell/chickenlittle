package team

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &teamEndpoint{}
	plugin.RegisterEndpoint(e)
}

type teamEndpoint struct {
	m *model.Model
	n *notification.Engine
}

func (a *teamEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/teams", a.List).Methods("GET")
	r.HandleFunc("/teams", a.Create).Methods("POST")
	r.HandleFunc("/teams/{team}", a.Show).Methods("GET")
	r.HandleFunc("/teams/{team}", a.Delete).Methods("DELETE")
	r.HandleFunc("/teams/{team}", a.Update).Methods("PUT")
	r.HandleFunc("/teams/{team}/notify", a.Notify).Methods("POST")
}

func (a *teamEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for teams
}

func (a *teamEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for teams
}

func (a *teamEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
}

func (a *teamEndpoint) Name() string {
	return "TeamEndpoint"
}
