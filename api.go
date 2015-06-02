package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chrissnell/victorops-go"
	"github.com/gorilla/mux"
)

type Response struct {
	People []Person `json:"people,omitempty"`
	Error  string
}

type Notification struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Priority victorops.MessageType `json:"priority,omitempty"`
}

func ListPeople(w http.ResponseWriter, r *http.Request) {
	var res Response

	for _, v := range config.People {
		res.People = append(res.People, *v.Sanitized())
	}

	json.NewEncoder(w).Encode(res)
}

func ShowPerson(w http.ResponseWriter, r *http.Request) {
	var res Response

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := config.GetPerson(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.People = append(res.People, *p.Sanitized())

	json.NewEncoder(w).Encode(res)
}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {

	var n []Notification
	var res Response

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*10))
	// If something went wrong, return an error in the JSON response
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = r.Body.Close()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = json.Unmarshal(body, &n)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	vo := victorops.NewClient(config.Integrations.VictorOps.APIKey)

	for _, person := range n {
		log.Println("Person:", person.User)
		p, err := config.GetPerson(person.User)
		if err != nil {
			res.Error = err.Error()
			log.Println("config.GetPerson() Error:", err)
			json.NewEncoder(w).Encode(res)
			return
		}

		e := &victorops.Event{
			RoutingKey:  p.VictorOpsRoutingKey,
			MessageType: victorops.Critical,
			EntityID:    person.Message,
			Timestamp:   time.Now(),
		}

		resp, err := vo.SendAlert(e)
		if err != nil {
			res.Error = err.Error()
			log.Println("Error:", err)
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Println("VO Response - Result:", resp.Result, "EntityID:", resp.EntityID, "Message:", resp.EntityID)
	}
	return

}
