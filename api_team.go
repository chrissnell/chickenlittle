package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TeamsResponse struct {
	Teams   []Team `json:"teams"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Fetches every team from the DB and returns them as JSON
func ListTeams(w http.ResponseWriter, r *http.Request) {
	var res TeamsResponse

	t, err := c.GetAllTeams()
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusInternalServerError)
		return
	}

	for _, v := range t {
		res.Teams = append(res.Teams, *v)
	}

	json.NewEncoder(w).Encode(res)
}

// Fetches a single team from the DB and returns them as JSON
func ShowTeam(w http.ResponseWriter, r *http.Request) {
	var res TeamsResponse

	vars := mux.Vars(r)
	name := vars["team"]

	p, err := c.GetTeam(name)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.Teams = append(res.Teams, *p)

	json.NewEncoder(w).Encode(res)
}

// Deletes the specified team from the database
func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	var res PeopleResponse

	vars := mux.Vars(r)
	name := vars["team"]

	err := c.DeleteTeam(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprint("Team ", name, " deleted")

	json.NewEncoder(w).Encode(res)
}

// Creates a new team in the database
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	var res TeamsResponse
	var t Team

	// We're getting the details of this new team from the POSTed JSON
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

	// Attempt to unmarshall the JSON into our Team struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// If a name was not provided, return an error
	if t.Name == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide a team name"
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure that this team doesn't already exist
	fp, err := c.GetTeam(t.Name)
	if fp != nil && fp.Name != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Team ", t.Name, " already exists. Use PUT /teams/", t.Name, "/ to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetTeam() failed:", err)
	}

	// Store our new team in the DB
	err = c.StoreTeam(&t)
	if err != nil {
		log.Println("Error storing team:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("Team ", t.Name, " created")

	json.NewEncoder(w).Encode(res)
}

// Updates an existing team in the database
func UpdateTeam(w http.ResponseWriter, r *http.Request) {
	var res TeamsResponse
	var t Team

	vars := mux.Vars(r)
	name := vars["team"]

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

	err = json.Unmarshal(body, &t)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure the team actually exists before updating
	fp, err := c.GetTeam(name)
	if fp != nil && fp.Name == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Team ", t.Name, " does not exist. Use POST to create.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetTeam() failed for", name)
	}

	// Now that we know our team exists in the DB, copy the name from the URI path and add it to our struct
	t.Name = name

	// Store the updated team in the DB
	err = c.StoreTeam(&t)
	if err != nil {
		log.Println("Error storing team", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Teams = append(res.Teams, t)
	res.Message = fmt.Sprint("Team ", name, " updated")

	json.NewEncoder(w).Encode(res)

}
