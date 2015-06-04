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

	p := Person{
		Username: "chris.snell",
		FullName: "Christopher Snell",
	}

	err = c.StorePerson(&p)
	if err != nil {
		log.Fatalf("Could not store person: %+v\n", p)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/people", ListPeople).
		Methods("GET")

	router.HandleFunc("/people/{person}/", ShowPerson).
		Methods("GET")

	router.HandleFunc("/people/{person}/", DeletePerson).
		Methods("DELETE")

	router.HandleFunc("/people/{person}", UpdatePerson).
		Methods("PUT")

	router.HandleFunc("/people", CreatePerson).
		Methods("POST")

	router.HandleFunc("/notify", NotifyPerson).
		Methods("POST")

	log.Fatal(http.ListenAndServe(c.Config.Service.ListenAddr, router))
}
