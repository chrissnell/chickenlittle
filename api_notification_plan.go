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

type NotificationPlanResponse struct {
	NotificationPlan NotificationPlan `json:"people,omitempty"`
	Message          string           `json:"message"`
	Error            string           `json:"error"`
}

func ShowNotificationPlan(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := c.GetNotificationPlan(username)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.NotificationPlan = *p

	json.NewEncoder(w).Encode(res)
}

func DeleteNotificationPlan(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse

	vars := mux.Vars(r)
	username := vars["person"]

	err := c.DeleteNotificationPlan(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprint("Notification plan for user ", username, " deleted")

	json.NewEncoder(w).Encode(res)
}

func CreateNotificationPlan(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse
	var p []NotificationStep

	vars := mux.Vars(r)
	username := vars["person"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*15))
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

	if username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username in URL"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := c.GetPerson(username)
	if fp != nil && fp.Username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", username, " does not exist. Create the user first before adding a notification plan for them.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	if len(p) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide at least one notification step in JSON"
		json.NewEncoder(w).Encode(res)
		return
	}

	np, err := c.GetNotificationPlan(username)
	if np != nil && np.Username != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Notification plan for user ", username, " already exists. Use PUT /plan/", username, " to update..")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetNotificationPlan() failed for", username)
	}

	plan := NotificationPlan{Username: username, Steps: p}

	err = c.StoreNotificationPlan(&plan)
	if err != nil {
		log.Println("Error storing notification plan:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("Notification plan for user ", username, " created")

	json.NewEncoder(w).Encode(res)
}

func UpdateNotificationPlan(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse
	var p []NotificationStep

	vars := mux.Vars(r)
	username := vars["person"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*15))
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

	if username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username in URL"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := c.GetPerson(username)
	if fp != nil && fp.Username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Can't update notification plan for user (", username, ") that does not exist.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	np, err := c.GetNotificationPlan(username)
	if (np != nil && np.Username == "") || err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Notification plan for user ", username, " doesn't exist. Use POST /plan/", username, " to create one first before attempting to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetNotificationPlan() failed for", username)
	}

	np.Steps = p

	err = c.StoreNotificationPlan(np)
	if err != nil {
		log.Println("Error storing notification plan:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.NotificationPlan = *np
	res.Message = fmt.Sprint("Notification plan for user ", username, " updated")

	json.NewEncoder(w).Encode(res)

}
