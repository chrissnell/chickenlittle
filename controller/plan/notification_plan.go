package plan

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

// ShowNotificationPlan returns a JSON-formatted NotificationPlan for a Person
func (a *npEndpoint) Show(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse

	vars := mux.Vars(r)
	username := vars["person"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	p, err := a.m.GetNotificationPlan(username)
	if err != nil {
		res.Error = err.Error()
		errjson, _ := json.Marshal(res)
		http.Error(w, string(errjson), http.StatusNotFound)
		return
	}

	res.NotificationPlan = *p

	json.NewEncoder(w).Encode(res)
}

// DeleteNotificationPlan deletes a Person's NotificationPlan
func (a *npEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse

	vars := mux.Vars(r)
	username := vars["person"]

	np, err := a.m.GetNotificationPlan(username)
	if np == nil {
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Notification plan for user ", username, " doesn't exist and thus, cannot be deleted")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetNotificationPlan() failed for", username)
	}

	err = a.m.DeleteNotificationPlan(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res.Message = fmt.Sprint("Notification plan for user ", username, " deleted")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(res)
}

// CreateNotificationPlan creates a NotificationPlan for a Person
func (a *npEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse
	var p model.NotificationPlan

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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

	if p.Username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username in URL"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := a.m.GetPerson(p.Username)
	if fp != nil && fp.Username == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("User ", p.Username, " does not exist. Create the user first before adding a notification plan for them.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	// The NotificationPlan provided must have at least one NotificationStep
	if len(p.Steps) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide at least one notification step in JSON"
		json.NewEncoder(w).Encode(res)
		return
	}

	np, err := a.m.GetNotificationPlan(p.Username)
	if np != nil && np.Username != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Notification plan for user ", p.Username, " already exists. Use PUT /plan/", p.Username, " to update..")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetNotificationPlan() failed for", p.Username)
	}

	err = a.m.StoreNotificationPlan(&p)
	if err != nil {
		log.Println("Error storing notification plan:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Message = fmt.Sprint("Notification plan for user ", p.Username, " created")

	json.NewEncoder(w).Encode(res)
}

// UpdateNotificationPlan updates an NotificationPlan for a Person
func (a *npEndpoint) Update(w http.ResponseWriter, r *http.Request) {
	var res NotificationPlanResponse
	var p model.NotificationPlan

	vars := mux.Vars(r)
	username := vars["person"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if username == "" {
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide username in URL"
		json.NewEncoder(w).Encode(res)
		return
	}

	fp, err := a.m.GetPerson(username)
	if fp != nil && fp.Username == "" {
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Can't update notification plan for user (", username, ") that does not exist.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetPerson() failed:", err)
	}

	// The NotificationPlan provided must have at least one NotificationStep
	if len(p.Steps) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Error = "Must provide at least one notification step in JSON"
		json.NewEncoder(w).Encode(res)
		return
	}

	np, err := a.m.GetNotificationPlan(username)
	if (np != nil && np.Username == "") || np == nil {
		w.WriteHeader(422) // unprocessable entity
		res.Error = fmt.Sprint("Notification plan for user ", username, " doesn't exist. Use POST /plan/", username, " to create one first before attempting to update.")
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.Println("GetNotificationPlan() failed for", username)
	}

	// Replace the NotificationSteps of the fetched plan with those from this request
	np.Steps = p.Steps

	err = a.m.StoreNotificationPlan(np)
	if err != nil {
		log.Println("Error storing notification plan:", err)
		w.WriteHeader(422) // unprocessable entity
		res.Error = err.Error()
		log.Println(res.Error)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.NotificationPlan = *np
	res.Message = fmt.Sprint("Notification plan for user ", username, " updated")

	json.NewEncoder(w).Encode(res)
}
