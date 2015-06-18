package controller

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

// API contains the HTTP API with all necessary dependencies
type Controller struct {
	c config.Config
	m *model.Model
	n *notification.Engine
}

// New will create a new API
func New(config config.Config, model *model.Model, eng *notification.Engine) *Controller {
	a := &Controller{
		c: config,
		m: model,
		n: eng,
	}
	return a
}

// APIRouter will create a new gorilla Router for handling all REST API calls
func (a *Controller) APIRouter() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true)

	apiRouter.HandleFunc("/people", a.ListPeople).Methods("GET")
	apiRouter.HandleFunc("/people", a.CreatePerson).Methods("POST")
	apiRouter.HandleFunc("/people/{person}", a.ShowPerson).Methods("GET")
	apiRouter.HandleFunc("/people/{person}", a.DeletePerson).Methods("DELETE")
	apiRouter.HandleFunc("/people/{person}", a.UpdatePerson).Methods("PUT")

	apiRouter.HandleFunc("/plan/{person}", a.CreateNotificationPlan).Methods("POST")
	apiRouter.HandleFunc("/plan/{person}", a.ShowNotificationPlan).Methods("GET")
	apiRouter.HandleFunc("/plan/{person}", a.DeleteNotificationPlan).Methods("DELETE")
	apiRouter.HandleFunc("/plan/{person}", a.UpdateNotificationPlan).Methods("PUT")

	apiRouter.HandleFunc("/people/{person}/notify", a.NotifyPerson).Methods("POST")
	apiRouter.HandleFunc("/notifications/{uuid}", a.StopNotification).Methods("DELETE")

	apiRouter.HandleFunc("/teams", a.ListTeams).Methods("GET")
	apiRouter.HandleFunc("/teams", a.CreateTeam).Methods("POST")
	apiRouter.HandleFunc("/teams/{team}", a.ShowTeam).Methods("GET")
	apiRouter.HandleFunc("/teams/{team}", a.DeleteTeam).Methods("DELETE")
	apiRouter.HandleFunc("/teams/{team}", a.UpdateTeam).Methods("PUT")

	return apiRouter
}

// CallbackRouter will create a new gorilla Router for handling all Callback Actions
func (a *Controller) CallbackRouter() *mux.Router {
	callbackRouter := mux.NewRouter().StrictSlash(true)

	callbackRouter.HandleFunc("/{uuid}/twiml/{action}", a.GenerateTwiML).Methods("POST")
	callbackRouter.HandleFunc("/{uuid}/callback", a.ReceiveCallback).Methods("POST")
	callbackRouter.HandleFunc("/{uuid}/digits", a.ReceiveDigits).Methods("POST")
	callbackRouter.HandleFunc("/sms", a.ReceiveSMSReply).Methods("POST")

	return callbackRouter
}

// ClickRouter will create a new gorilla Router for handling all Click Actions
func (a *Controller) ClickRouter() *mux.Router {
	clickRouter := mux.NewRouter().StrictSlash(true)

	clickRouter.HandleFunc("/{uuid}/stop", a.StopNotificationClick).Methods("GET")

	return clickRouter
}
