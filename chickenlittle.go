package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"time"

	"github.com/gorilla/mux"

	"gopkg.in/yaml.v2"
)

var (
	c ChickenLittle
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

	np := &NotificationPlan{}
	np.Username = "evan.snell"
	np.Steps = append(np.Steps, NotificationStep{Method: Voice,
		Data:              "2108593107",
		NotifyUntilPeriod: time.Minute * 5,
		NotifyEveryPeriod: time.Minute,
	})

	err = c.StoreNotificationPlan(np)
	if err != nil {
		log.Fatalln("Could not store notification plan:", np)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/people", ListPeople).
		Methods("GET")

	router.HandleFunc("/people", CreatePerson).
		Methods("POST")

	router.HandleFunc("/people/{person}/", ShowPerson).
		Methods("GET")

	router.HandleFunc("/people/{person}/", DeletePerson).
		Methods("DELETE")

	router.HandleFunc("/people/{person}", UpdatePerson).
		Methods("PUT")

	router.HandleFunc("/plan", CreateNotificationPlan).
		Methods("POST")

	router.HandleFunc("/plan/{person}/", ShowNotificationPlan).
		Methods("GET")

	router.HandleFunc("/plan/{person}/", DeleteNotificationPlan).
		Methods("DELETE")

	router.HandleFunc("/plan/{person}", UpdateNotificationPlan).
		Methods("PUT")

	router.HandleFunc("/people/{person}/notify", NotifyPerson).
		Methods("POST")

	log.Fatal(http.ListenAndServe(c.Config.Service.ListenAddr, router))
}
