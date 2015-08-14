package notification

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &notificationEndpoint{}
	plugin.RegisterEndpoint(e)
}

type notificationEndpoint struct {
	m *model.Model
	n *notification.Engine
}

func (a *notificationEndpoint) APIRoutes(r *mux.Router) {
	r.HandleFunc("/notifications/{uuid}", a.StopNotification).Methods("DELETE")
}

func (a *notificationEndpoint) CallbackRoutes(r *mux.Router) {
	// no callback routes for notification endpoint
}

func (a *notificationEndpoint) ClickRoutes(r *mux.Router) {
	r.HandleFunc("/{uuid}/stop", a.StopNotificationClick).Methods("GET")
}

func (a *notificationEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.m = m
	a.n = n
}

func (a *notificationEndpoint) Name() string {
	return "NotificationEndpoint"
}
