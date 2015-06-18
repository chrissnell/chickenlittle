package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/chrissnell/chickenlittle/ne"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile *string
	c       ChickenLittle
)

type ChickenLittle struct {
	Config Config
	DB     DB
	Notify *ne.Engine
}

func main() {

	cfgFile = flag.String("config", "config.yaml", "Path to config file (default: ./config.yaml)")
	flag.Parse()

	// Read our server configuration
	filename, _ := filepath.Abs(*cfgFile)
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error opening config file.  Did you pass the -config flag?  Run with -h for help.\n", err)
	}
	err = yaml.Unmarshal(cfgFile, &c.Config)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	// Open our BoltDB handle
	c.DB.Open(c.Config.Service.DBFile)
	defer c.DB.Close()

	// Launch the notification engine
	neConfig := ne.Config{}
	neConfig.Twilio.AccountSID = c.Config.Integrations.Twilio.AccountSID
	neConfig.Twilio.AuthToken = c.Config.Integrations.Twilio.AuthToken
	neConfig.Twilio.CallFromNumber = c.Config.Integrations.Twilio.CallFromNumber
	neConfig.Twilio.APIBaseURL = c.Config.Integrations.Twilio.APIBaseURL
	neConfig.Mailgun.Enabled = c.Config.Integrations.Mailgun.Enabled
	neConfig.Mailgun.APIKey = c.Config.Integrations.Mailgun.APIKey
	neConfig.Mailgun.Hostname = c.Config.Integrations.Mailgun.Hostname
	neConfig.SMTP.Hostname = c.Config.Integrations.SMTP.Hostname
	neConfig.SMTP.Port = c.Config.Integrations.SMTP.Port
	neConfig.SMTP.Login = c.Config.Integrations.SMTP.Login
	neConfig.SMTP.Password = c.Config.Integrations.SMTP.Password
	neConfig.SMTP.Sender = c.Config.Integrations.SMTP.Sender
	neConfig.Service.ClickURLBase = c.Config.Service.ClickURLBase
	neConfig.Service.CallbackURLBase = c.Config.Service.CallbackURLBase
	c.Notify = ne.New(neConfig)

	// Set up our API endpoint router
	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.APIListenAddr, apiRouter()))
	}()

	// Set up our Twilio callback endpoint router
	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.CallbackListenAddr, callbackRouter()))
	}()

	// Set up our Click endpoint router to handle stop requests from browsers
	log.Fatal(http.ListenAndServe(c.Config.Service.ClickListenAddr, clickRouter()))
}

func apiRouter() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true)

	apiRouter.HandleFunc("/people", ListPeople).
		Methods("GET")

	apiRouter.HandleFunc("/people", CreatePerson).
		Methods("POST")

	apiRouter.HandleFunc("/people/{person}", ShowPerson).
		Methods("GET")

	apiRouter.HandleFunc("/people/{person}", DeletePerson).
		Methods("DELETE")

	apiRouter.HandleFunc("/people/{person}", UpdatePerson).
		Methods("PUT")

	apiRouter.HandleFunc("/plan/{person}", CreateNotificationPlan).
		Methods("POST")

	apiRouter.HandleFunc("/plan/{person}", ShowNotificationPlan).
		Methods("GET")

	apiRouter.HandleFunc("/plan/{person}", DeleteNotificationPlan).
		Methods("DELETE")

	apiRouter.HandleFunc("/plan/{person}", UpdateNotificationPlan).
		Methods("PUT")

	apiRouter.HandleFunc("/people/{person}/notify", NotifyPerson).
		Methods("POST")

	apiRouter.HandleFunc("/notifications/{uuid}", StopNotification).
		Methods("DELETE")

	apiRouter.HandleFunc("/teams", ListTeams).
		Methods("GET")

	apiRouter.HandleFunc("/teams", CreateTeam).
		Methods("POST")

	apiRouter.HandleFunc("/teams/{team}", ShowTeam).
		Methods("GET")

	apiRouter.HandleFunc("/teams/{team}", DeleteTeam).
		Methods("DELETE")

	apiRouter.HandleFunc("/teams/{team}", UpdateTeam).
		Methods("PUT")

	return apiRouter
}

func callbackRouter() *mux.Router {
	callbackRouter := mux.NewRouter().StrictSlash(true)

	callbackRouter.HandleFunc("/{uuid}/twiml/{action}", GenerateTwiML).
		Methods("POST")

	callbackRouter.HandleFunc("/{uuid}/callback", ReceiveCallback).
		Methods("POST")

	callbackRouter.HandleFunc("/{uuid}/digits", ReceiveDigits).
		Methods("POST")

	callbackRouter.HandleFunc("/sms", ReceiveSMSReply).
		Methods("POST")

	return callbackRouter
}

func clickRouter() *mux.Router {
	clickRouter := mux.NewRouter().StrictSlash(true)

	clickRouter.HandleFunc("/{uuid}/stop", StopNotificationClick).
		Methods("GET")

	return clickRouter
}
