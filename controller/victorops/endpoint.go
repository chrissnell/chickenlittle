package victorops

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &victoropsEndpoint{}
	plugin.RegisterEndpoint(e)
}

type victoropsEndpoint struct {
	c config.Config
	m *model.Model
}

func (a *victoropsEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/person/{person}/notifyvo", a.Notify).Methods("POST")
}

func (a *victoropsEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for victorops
}

func (a *victoropsEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for victorops
}

func (a *victoropsEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.c = c
	a.m = m
}
func (a *victoropsEndpoint) Name() string {
	return "VictorOpsEndpoint"
}
