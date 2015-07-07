package controller

import (
	"log"

	"github.com/chrissnell/chickenlittle/config"
	_ "github.com/chrissnell/chickenlittle/controller/escalation"   // register escalation plugin
	_ "github.com/chrissnell/chickenlittle/controller/notification" // register notification plugin
	_ "github.com/chrissnell/chickenlittle/controller/people"       // register people plugin
	_ "github.com/chrissnell/chickenlittle/controller/plan"         // register escalation plan plugin
	"github.com/chrissnell/chickenlittle/controller/plugin"
	_ "github.com/chrissnell/chickenlittle/controller/rotation"  // register rotation plugin
	_ "github.com/chrissnell/chickenlittle/controller/team"      // register team plugin
	_ "github.com/chrissnell/chickenlittle/controller/twilio"    // register twilio plugin
	_ "github.com/chrissnell/chickenlittle/controller/victorops" // register victorops plugin
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/chrissnell/chickenlittle/rotation"
	"github.com/gorilla/mux"
)

// Controller contains the HTTP controller with all necessary dependencies
type Controller struct {
	c   config.Config
	m   *model.Model
	n   *notification.Engine
	r   *rotation.Engine
	eps []plugin.Endpoint
}

// New will create a new Controller
func New(config config.Config, model *model.Model, eng *notification.Engine, rot *rotation.Engine) *Controller {
	a := &Controller{
		c:   config,
		m:   model,
		n:   eng,
		r:   rot,
		eps: make([]plugin.Endpoint, 0, 10),
	}
	for _, ep := range plugin.Endpoints() {
		ep.SetConfig(model, config, eng)
		a.eps = append(a.eps, ep)
		log.Println("Plugin configured:", ep.Name())
	}
	return a
}

// APIRouter will create a new gorilla Router for handling all REST API calls
func (a *Controller) APIRouter() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true)

	for _, ep := range a.eps {
		ep.APIRoutes(apiRouter)
	}

	return apiRouter
}

// CallbackRouter will create a new gorilla Router for handling all Callback Actions
func (a *Controller) CallbackRouter() *mux.Router {
	callbackRouter := mux.NewRouter().StrictSlash(true)

	for _, ep := range a.eps {
		ep.CallbackRoutes(callbackRouter)
	}

	return callbackRouter
}

// ClickRouter will create a new gorilla Router for handling all Click Actions
func (a *Controller) ClickRouter() *mux.Router {
	clickRouter := mux.NewRouter().StrictSlash(true)

	for _, ep := range a.eps {
		ep.ClickRoutes(clickRouter)
	}

	return clickRouter
}
