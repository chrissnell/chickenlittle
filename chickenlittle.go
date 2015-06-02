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
	config Config
)

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	cfgFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalln("Error:", err)
	}

	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/people", ListPeople).
		Methods("GET")

	router.HandleFunc("/people/{person}/", ShowPerson).
		Methods("GET")

	router.HandleFunc("/notify", NotifyPerson).
		Methods("POST")

	log.Fatal(http.ListenAndServe(config.Service.ListenAddr, router))
}
