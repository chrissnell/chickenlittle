package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/victorops-go"
	"github.com/gorilla/mux"
)

// PeopleResponse is the JSON response
type PeopleResponse struct {
	People  []model.Person `json:"people"`
	Message string         `json:"message"`
	Error   string         `json:"error"`
}

// Notification is the struct of an incoming notification request
type Notification struct {
	Username string                `json:"username"`
	Message  string                `json:"message"`
	Priority victorops.MessageType `json:"priority,omitempty"`
}

// ListPeople fetches every person from the DB and returns them as JSON
func (a *Controller) ListPeople(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	p, err := a.m.GetAllPeople()
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

// ShowPerson fetches a single person form the DB and returns them as JSON
func (a *Controller) ShowPerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := a.m.GetPerson(username)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.People = append(res.People, *p)

	json.NewEncoder(w).Encode(res)
}

// DeletePerson deletes the specified person from the database
func (a *Controller) DeletePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	vars := mux.Vars(r)
	username := vars["person"]

	// Make sure the user actually exists before updating
	p, err := a.m.GetPerson(username)
	if p == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", username, " does not exist and thus, cannot be deleted")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed for", username)
	}

	err = a.m.DeletePerson(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprint("User ", username, " deleted")

	json.NewEncoder(w).Encode(res)
}

// CreatePerson creates a new person in the database
func (a *Controller) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse
	var p model.Person

	// We're getting the details of this new person from the POSTed JSON
	// so we first need to read in the body of the POST
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

	// Attempt to unmarshall the JSON into our Person struct
	err = json.Unmarshal(body, &p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// If a username *and* fullname were not provided, return an error
	if p.Username == "" || p.FullName == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username and fullname"
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure that this user doesn't already exist
	fp, err := a.m.GetPerson(p.Username)
	if fp != nil && fp.Username != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Username, " already exists. Use PUT /people/", p.Username, " to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	// Store our new person in the DB
	err = a.m.StorePerson(&p)
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

// UpdatePerson updates an existing person in the database
func (a *Controller) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse
	var p model.Person

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

	// The only field that can be updated currently is fullname, so make sure one was provided
	if p.FullName == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide a fullname to update"
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure the user actually exists before updating
	fp, err := a.m.GetPerson(username)
	if (fp != nil && fp.Username == "") || fp == nil {
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

	// Store the updated user in the DB
	err = a.m.StorePerson(&p)
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
