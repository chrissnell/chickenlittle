package escalation

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chrissnell/chickenlittle/model"
	"github.com/gorilla/mux"
)

func (a *escalationEndpoint) Show(w http.ResponseWriter, r *http.Request) {
	var res Response

	vars := mux.Vars(r)
	name := vars["plan"]

	p, err := a.m.GetEscalationPlan(name)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.Plans = append(res.Plans, *p)

	json.NewEncoder(w).Encode(res)
}

func (a *escalationEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	var res Response

	vars := mux.Vars(r)
	name := vars["plan"]

	p, err := a.m.GetEscalationPlan(name)
	if p == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprintf("Escalation Plan %s does not exists and thus, cannot be deleted", name)
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetEscalationPlan() failed for", name)
	}

	err = a.m.DeleteEscalationPlan(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprintf("Escalation Plan %s deleted", name)

	json.NewEncoder(w).Encode(res)
}

func (a *escalationEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var res Response
	var p model.EscalationPlan

	// We're getting the details of this new plan from the POSTed JSON
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

	// Attempt to unmarshall the JSON into our EscalationPlan struct
	err = json.Unmarshal(body, &p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// If a name was not provided, return an error
	if p.Name == "" { // TODO further checks
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide a name"
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure that this plan doesn't already exist
	fp, err := a.m.GetEscalationPlan(p.Name)
	if fp != nil && fp.Name != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Escalation Plan ", p.Name, " already exists. Use PUT /plan/", p.Name, " to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetEscalationPlan() failed:", err)
	}

	// Store our new plan in the DB
	err = a.m.StoreEscalationPlan(&p)
	if err != nil {
		log.Println("Error storing plan:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("Pla ", p.Name, " created")

	json.NewEncoder(w).Encode(res)
}

func (a *escalationEndpoint) Update(w http.ResponseWriter, r *http.Request) {
	var res Response
	var p model.EscalationPlan

	vars := mux.Vars(r)
	name := vars["plan"]

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

	// TODO validate updated fields
	if p.Name == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide a name to update"
		json.NewEncoder(w).Encode(res)
		return
	}

	// Make sure the plan actually exists before updating
	fp, err := a.m.GetEscalationPlan(name)
	if (fp != nil && fp.Name == "") || fp == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Name, " does not exist. Use POST to create.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetEscalationPlan() failed for", name)
	}

	// Now that we know our plan exists in the DB, copy the name from the URI path and add it to our struct
	p.Name = name

	// Store the updated user in the DB
	err = a.m.StoreEscalationPlan(&p)
	if err != nil {
		log.Println("Error storing person:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Plans = append(res.Plans, p)
	res.Message = fmt.Sprint("Escalation Plan ", name, " updated")

	json.NewEncoder(w).Encode(res)

}
