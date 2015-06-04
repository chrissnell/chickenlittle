package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chrissnell/victorops-go"
	"github.com/gorilla/mux"
)

type PeopleResponse struct {
	People  []Person `json:"people,omitempty"`
	Message string   `json:"message"`
	Error   string   `json:"error"`
}

type Notification struct {
	Username string                `json:"username"`
	Message  string                `json:"message"`
	Priority victorops.MessageType `json:"priority,omitempty"`
}

func ListPeople(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	p, err := c.GetAllPeople()
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusInternalServerError)
		return
	}

	for _, v := range p {
		res.People = append(res.People, *v)
	}

	json.NewEncoder(w).Encode(res)
}

func ShowPerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := c.GetPerson(username)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.People = append(res.People, *p.Sanitized())

	json.NewEncoder(w).Encode(res)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	vars := mux.Vars(r)
	username := vars["person"]

	err := c.DeletePerson(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprint("User ", username, " deleted")

	json.NewEncoder(w).Encode(res)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse
	var p Person

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

	err = json.Unmarshal(body, &p)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if p.Username == "" || p.FullName == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username and fullname"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := c.GetPerson(p.Username)
	if fp != nil && fp.Username != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Username, " already exists. Use PUT /people/", p.Username, "/ to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	err = c.StorePerson(&p)
	if err != nil {
		log.Println("Error storing person:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("User ", p.Username, " created")

	json.NewEncoder(w).Encode(res)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse
	var p Person

	vars := mux.Vars(r)
	username := vars["person"]

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

	err = json.Unmarshal(body, &p)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if p.FullName == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide a fullname to update"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := c.GetPerson(username)
	if fp != nil && fp.Username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Username, " does not exist. Use POST to create.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed for", username)
	}

	// Now that we know our user exists in the DB, copy the username from the URI path and add it to our struct
	p.Username = username

	err = c.StorePerson(&p)
	if err != nil {
		log.Println("Error storing person:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.People = append(res.People, p)
	res.Message = fmt.Sprint("User ", username, " updated")

	json.NewEncoder(w).Encode(res)

}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {

	var n Notification
	var res PeopleResponse

	vars := mux.Vars(r)
	username := vars["person"]

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

	vo := victorops.NewClient(c.Config.Integrations.VictorOps.APIKey)

	p, err := c.GetPerson(username)
	if err != nil {
		res.Error = err.Error()
		log.Println("config.GetPerson() Error:", err)
		json.NewEncoder(w).Encode(res)
		return
	}

	e := &victorops.Event{
		RoutingKey:  p.VictorOpsRoutingKey,
		MessageType: victorops.Critical,
		EntityID:    n.Message,
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
	return

}
