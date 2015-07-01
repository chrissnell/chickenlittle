package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/controller"
	"github.com/chrissnell/chickenlittle/db"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
)

// ChickenLittle contains the notification service server
type ChickenLittle struct {
	Config config.Config
	db     *db.DB
	notify *notification.Engine
	api    *controller.Controller
	model  *model.Model
}

// New creates a new instance of ChickenLittle with the given configuration file
func New(filename string) *ChickenLittle {
	c := &ChickenLittle{}

	// Read our server configuration
	filename, _ = filepath.Abs(filename)
	cfg, err := config.New(filename)
	if err != nil {
		log.Fatalln("Error reading config file.  Did you pass the -config flag?  Run with -h for help.\n", err)
	}
	c.Config = cfg

	// Open our BoltDB handle
	c.db = db.New(c.Config.Service.DBFile)
	defer c.db.Close()

	// Initialize the data model
	c.model = model.New(c.db)

	// Launch the notification engine
	c.notify = notification.New(c.Config)

	// Initialize the Controller
	c.api = controller.New(c.Config, c.model, c.notify)

	return c
}

// Listen will start the HTTP listeners for API, Click and Callback Routers.
func (c *ChickenLittle) Listen() {
	// Set up our API endpoint router
	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.APIListenAddr, c.api.APIRouter()))
	}()

	// Set up our Twilio callback endpoint router
	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.CallbackListenAddr, c.api.CallbackRouter()))
	}()

	// Set up our Click endpoint router to handle stop requests from browsers
	log.Fatal(http.ListenAndServe(c.Config.Service.ClickListenAddr, c.api.ClickRouter()))
}

// Close will shut down the ChickenLittle service
func (c *ChickenLittle) Close() {
	c.db.Close()
}

func main() {
	cfgFile := flag.String("config", "config.yaml", "Path to config file (default: ./config.yaml)")
	flag.Parse()

	c := New(*cfgFile)
	defer c.Close()
	c.Listen()
}
