package plugin

import (
	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/gorilla/mux"
)

// Endpoint defines the common interface for the different controller endpoints
type Endpoint interface {
	APIRoutes(*mux.Router)
	CallbackRoutes(*mux.Router)
	ClickRoutes(*mux.Router)
	SetConfig(*model.Model, config.Config, *notification.Engine)
	Name() string
}
