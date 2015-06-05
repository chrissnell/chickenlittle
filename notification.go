package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

var (
	stopChan   chan string
	timerChan  <-chan time.Time
	tickerChan <-chan time.Time
)

type NotificationsInProgress struct {
	Stoppers map[string]chan bool
	Mu       sync.Mutex
}

type NotifyPersonResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

func StartNotificationEngine() {
	var NIP NotificationsInProgress

	// Initialize our map of Stopper channels
	// UUID -> channel
	NIP.Stoppers = make(map[string]chan bool)

	log.Println("StartNotificationEngine()")

	for {

		select {
		case np := <-planChan:
			// We've received a new notification plan
			log.Printf("Got plan: %+v\n", np)

			// Get the plan's UUID
			id := np.ID.String()

			NIP.Mu.Lock()

			// Create a new Stopper channel for this plan
			NIP.Stoppers[id] = make(chan bool)

			// Launch a goroutine to handle plan processing
			go planHandler(np, NIP.Stoppers[id])

			NIP.Mu.Unlock()
		case stopUUID := <-stopChan:
			// We've received a request to stop a notification plan
			NIP.Mu.Lock()

			// Check to see if the requested UUID is actually in progress
			_, prs := NIP.Stoppers[stopUUID]
			if prs {
				// It's in progress, so we'll send a message on its Stopper to
				// be received by the goroutine executing the plan
				NIP.Stoppers[stopUUID] <- true
			}
			NIP.Mu.Unlock()
		}
	}

}

func StopNotification(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	id := vars["uuid"]

	// Attempt to stop the notification by sending the UUID to the notification engine
	stopChan <- id

	// TO DO: make sure that this is a valid UUID and obtain
	//        confirmation of deletion

	res = NotifyPersonResponse{
		Message: "Attempting to terminate notification",
		UUID:    id,
	}

	json.NewEncoder(w).Encode(res)
}

func planHandler(np *NotificationPlan, sc <-chan bool) {
	log.Println("planHandler()")

	uuid := np.ID.String()

	for n, s := range np.Steps {

		log.Println("[", uuid, "]", "STEP", n)
		log.Println("[", uuid, "]", "Method:", s.Method)

		if n == len(np.Steps)-1 {
			// Last step, so we use a Ticker and NotifyEveryPeriod
			tickerChan = time.NewTicker(s.NotifyEveryPeriod).C
			log.Println("[", uuid, "]", "[Waiting", strconv.FormatFloat(s.NotifyEveryPeriod.Minutes(), 'f', 1, 64), "minutes]")

		} else {
			// Not the last step, so we use a Timer and NotifyUntilPeriod
			timerChan = time.NewTimer(s.NotifyUntilPeriod).C
			log.Println("[", uuid, "]", "[Waiting", strconv.FormatFloat(s.NotifyUntilPeriod.Minutes(), 'f', 1, 64), "minutes]")
		}

	timerLoop:
		for {
			select {
			case <-timerChan:
				log.Println("[", uuid, "]", "Step timer expired.  Proceeding!")
				break timerLoop
			case <-tickerChan:
				log.Println("[", uuid, "]", "**Tick**  Retry contact method!")
				log.Println("[", uuid, "]", "Waiting", strconv.FormatFloat(s.NotifyEveryPeriod.Minutes(), 'f', 1, 64), "minutes]")
			case <-sc:
				log.Println("[", uuid, "]", "Manual stop received.  Exiting loops.")
				return
			}
		}

		log.Println("Loop broken")
	}
}

func NotifyPerson(w http.ResponseWriter, r *http.Request) {
	var res NotifyPersonResponse

	vars := mux.Vars(r)
	username := vars["person"]

	p, err := c.GetNotificationPlan(username)
	if err != nil {
		// res.Error = err.Error()
		// errjson, _ := json.Marshal(res)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assign a UUID
	uuid.SwitchFormat(uuid.CleanHyphen)
	p.ID = uuid.NewV4()

	planChan <- p

	res = NotifyPersonResponse{
		Message:  "Notification initiated",
		UUID:     p.ID.String(),
		Username: p.Username,
	}

	json.NewEncoder(w).Encode(res)

}
