package twilio

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller/plugin"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

func init() {
	e := &twilioEndpoint{}
	plugin.RegisterEndpoint(e)
}

type twilioEndpoint struct {
	c config.Config
	n *notification.Engine
}

func (a *twilioEndpoint) APIRoutes(r *mux.Router) {
	// no api routes for twilio
}

func (a *twilioEndpoint) CallbackRoutes(r *mux.Router) {
	r.HandleFunc("/{uuid}/twiml/{action}", a.GenerateTwiML).Methods("POST")
	r.HandleFunc("/{uuid}/callback", a.ReceiveCallback).Methods("POST")
	r.HandleFunc("/{uuid}/digits", a.ReceiveDigits).Methods("POST")
	r.HandleFunc("/sms", a.ReceiveSMSReply).Methods("POST")
}

func (a *twilioEndpoint) ClickRoutes(r *mux.Router) {
	// no click routes for twilio
}

func (a *twilioEndpoint) SetConfig(m *model.Model, c config.Config, n *notification.Engine) {
	a.c = c
	a.n = n
}

func (a *twilioEndpoint) Name() string {
	return "TwilioEndpoint"
}
