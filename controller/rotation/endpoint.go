package rotation

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &rotationEndpoint{}
	plugin.RegisterEndpoint(e)
}

type rotationEndpoint struct {
	m *model.Model
}

func (a *rotationEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/rotation", a.Create).Methods("POST")
	r.HandleFunc("/rotation/{policy}", a.Show).Methods("GET")
	r.HandleFunc("/rotation/{policy}", a.Delete).Methods("DELETE")
	r.HandleFunc("/rotation/{policy}", a.Update).Methods("PUT")
}

func (a *rotationEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for rotation plan
}

func (a *rotationEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for rotation plan
}

func (a *rotationEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
}

func (a *rotationEndpoint) Name() string {
	return "RotationPolicyEndpoint"
}
