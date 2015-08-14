package victorops

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

// NotifyPersonViaVictorops will send a notification via VictorOps
func (a *victoropsEndpoint) Notify(w http.ResponseWriter, r *http.Request) {

	var n NotificationRequest
	var res NotificationResponse

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

	vo := victorops.NewClient(a.c.Integrations.VictorOps.APIKey)

	p, err := a.m.GetPerson(username)
	if err != nil {
		res.Error = err.Error()
		log.Println("config.GetPerson() Error:", err)
		json.NewEncoder(w).Encode(res)
		return
	}

	e := &victorops.Event{
		RoutingKey:  p.VictorOpsRoutingKey,
		MessageType: victorops.Critical,
		EntityID:    n.Content,
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
