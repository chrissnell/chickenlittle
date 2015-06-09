package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var (
	c        ChickenLittle
	NIP      NotificationsInProgress
	planChan = make(chan *NotificationRequest)
)

type ChickenLittle struct {
	Config Config
	DB     DB
}

func main() {

	// Read our server configuration
	filename, _ := filepath.Abs("./config.yaml")
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	err = yaml.Unmarshal(cfgFile, &c.Config)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	// Open our BoltDB handle
	c.DB.Open(c.Config.Service.DBFile)
	defer c.DB.Close()

	// Creat our stop channel and launch the notification engine
	stopChan = make(chan string)
	go StartNotificationEngine()

	// Set up our API endpoint router
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

	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.APIListenAddr, apiRouter))
	}()

	// Set up our Twilio callback endpoint router
	callbackRouter := mux.NewRouter().StrictSlash(true)

	callbackRouter.HandleFunc("/{uuid}/twiml/{action}", GenerateTwiML).
		Methods("POST")

	callbackRouter.HandleFunc("/{uuid}/callback", ReceiveCallback).
		Methods("POST")

	callbackRouter.HandleFunc("/{uuid}/digits", ReceiveDigits).
		Methods("POST")

	callbackRouter.HandleFunc("/sms", ReceiveSMSReply).
		Methods("POST")

	go func() {
		log.Fatal(http.ListenAndServe(c.Config.Service.CallbackListenAddr, callbackRouter))
	}()

	// Set up our Click endpoint router to handle stop requests from browsers
	clickRouter := mux.NewRouter().StrictSlash(true)

	clickRouter.HandleFunc("/{uuid}/stop", StopNotificationClick).
		Methods("GET")

	log.Fatal(http.ListenAndServe(c.Config.Service.ClickListenAddr, clickRouter))

}
