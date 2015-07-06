package controller

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

type RotationPolicyResponse struct {
	Policies []model.RotationPolicy `json:"policies"`
	Message  string                 `json:"message"`
	Error    string                 `json:"error"`
}

func (a *Controller) ShowRotationPolicy(w http.ResponseWriter, r *http.Request) {
	var res RotationPolicyResponse

	vars := mux.Vars(r)
	name := vars["policy"]

	p, err := a.m.GetRotationPolicy(name)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.Policies = append(res.Policies, *p)

	json.NewEncoder(w).Encode(res)
}

func (a *Controller) DeleteRotationPolicy(w http.ResponseWriter, r *http.Request) {
	var res RotationPolicyResponse

	vars := mux.Vars(r)
	name := vars["policy"]

	p, err := a.m.GetRotationPolicy(name)
	if p == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprintf("Escalation policy %s does not exists and thus, cannot be deleted", name)
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetRotationPolicy() failed for", name)
	}

	err = a.m.DeleteRotationPolicy(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprintf("Escalation policy %s deleted", name)

	json.NewEncoder(w).Encode(res)
}

func (a *Controller) CreateRotationPolicy(w http.ResponseWriter, r *http.Request) {
	var res RotationPolicyResponse
	var p model.RotationPolicy

	// We're getting the details of this new policy from the POSTed JSON
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

	// Attempt to unmarshall the JSON into our RotationPolicy struct
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

	// Make sure that this policy doesn't already exist
	fp, err := a.m.GetRotationPolicy(p.Name)
	if fp != nil && fp.Name != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Escalation policy ", p.Name, " already exists. Use PUT /rotation/", p.Name, " to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetRotationPolicy() failed:", err)
	}

	// Store our new policy in the DB
	err = a.m.StoreRotationPolicy(&p)
	if err != nil {
		log.Println("Error storing policy:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("Plan ", p.Name, " created")

	json.NewEncoder(w).Encode(res)
}

func (a *Controller) UpdateRotationPolicy(w http.ResponseWriter, r *http.Request) {
	var res RotationPolicyResponse
	var p model.RotationPolicy

	vars := mux.Vars(r)
	name := vars["policy"]

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

	// Make sure the policy actually exists before updating
	fp, err := a.m.GetRotationPolicy(name)
	if (fp != nil && fp.Name == "") || fp == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Name, " does not exist. Use POST to create.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetRotationPolicy() failed for", name)
	}

	// Now that we know our policy exists in the DB, copy the name from the URI path and add it to our struct
	p.Name = name

	// Store the updated user in the DB
	err = a.m.StoreRotationPolicy(&p)
	if err != nil {
		log.Println("Error storing person:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Policies = append(res.Policies, p)
	res.Message = fmt.Sprint("Escalation policy ", name, " updated")

	json.NewEncoder(w).Encode(res)

}
