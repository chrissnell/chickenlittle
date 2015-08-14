package people

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &personEndpoint{}
	plugin.RegisterEndpoint(e)
}

type personEndpoint struct {
	m *model.Model
	n *notification.Engine
}

func (a *personEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/people", a.List).Methods("GET")
	r.HandleFunc("/people", a.Create).Methods("POST")
	r.HandleFunc("/people/{person}", a.Show).Methods("GET")
	r.HandleFunc("/people/{person}", a.Delete).Methods("DELETE")
	r.HandleFunc("/people/{person}", a.Update).Methods("PUT")
	r.HandleFunc("/people/{person}/notify", a.Notify).Methods("POST")
}

func (a *personEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for people
}

func (a *personEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for people
}

func (a *personEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
	a.n = n
}

func (a *personEndpoint) Name() string {
	return "PersonEndpoint"
}
