package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func StartNotificationEngine() {

	log.Println("StartNotificationEngine()")

	for {
		np := <-planChan
		log.Printf("Got plan: %+v\n", np)
		planHandler(np)
	}

}

func planHandler(np *NotificationPlan) {
	log.Println("planHandler()")

	for _, s := range np.Steps {
		log.Println("Method:", s.Method, "Data:", s.Data, "Until", strconv.FormatFloat(s.NotifyUntilPeriod.Minutes(), 'f', 1, 64))
	}

}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["person"]

	p, err := c.GetNotificationPlan(username)
	if err != nil {
		// res.Error = err.Error()
		// errjson, _ := json.Marshal(res)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	planChan <- p

}
