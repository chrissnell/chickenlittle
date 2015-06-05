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
	planChan = make(chan *NotificationPlan)
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

	stopChan = make(chan string)
	go StartNotificationEngine()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/people", ListPeople).
		Methods("GET")

	router.HandleFunc("/people", CreatePerson).
		Methods("POST")

	router.HandleFunc("/people/{person}", ShowPerson).
		Methods("GET")

	router.HandleFunc("/people/{person}", DeletePerson).
		Methods("DELETE")

	router.HandleFunc("/people/{person}", UpdatePerson).
		Methods("PUT")

	router.HandleFunc("/plan/{person}", CreateNotificationPlan).
		Methods("POST")

	router.HandleFunc("/plan/{person}", ShowNotificationPlan).
		Methods("GET")

	router.HandleFunc("/plan/{person}", DeleteNotificationPlan).
		Methods("DELETE")

	router.HandleFunc("/plan/{person}", UpdateNotificationPlan).
		Methods("PUT")

	router.HandleFunc("/people/{person}/notify", NotifyPerson).
		Methods("POST")

	router.HandleFunc("/notifications/{uuid}", StopNotification).
		Methods("DELETE")

	log.Fatal(http.ListenAndServe(c.Config.Service.ListenAddr, router))
}
