package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	stopChan chan string
)

type NotificationsInProgress struct {
	Stoppers      map[string]chan bool
	Messages      map[string]string
	Conversations map[string]string
	Mu            sync.Mutex
}

func StartNotificationEngine() {
	// Initialize our map of Stopper channels
	// UUID -> channel
	NIP.Stoppers = make(map[string]chan bool)

	// Initialize our map of Messages
	NIP.Messages = make(map[string]string)

	// Initialize our map of Conversations
	NIP.Conversations = make(map[string]string)

	log.Println("StartNotificationEngine()")

	for {

		select {
		// IMPROVEMENT: We could implement a close() of planChan to indicate that the service is shutting down
		//              and instruct all notifications to cease
		case nr := <-planChan:
			// We've received a new notification plan

			// Get the plan's UUID
			id := nr.Plan.ID.String()

			NIP.Mu.Lock()

			// Create a new Stopper channel for this plan
			NIP.Stoppers[id] = make(chan bool)

			// Save the message to NIP.Message
			NIP.Messages[id] = nr.Content

			// Launch a goroutine to handle plan processing
			go notificationHandler(nr, NIP.Stoppers[id])

			NIP.Mu.Unlock()
		case stopUUID := <-stopChan:
			// We've received a request to stop a notification plan
			NIP.Mu.Lock()

			// Check to see if the requested UUID is actually in progress
			_, prs := NIP.Stoppers[stopUUID]
			if prs {

				log.Println("[", stopUUID, "]", "Sending a stop notification to the plan processor")

				// It's in progress, so we'll send a message on its Stopper to
				// be received by the goroutine executing the plan
				NIP.Stoppers[stopUUID] <- true
			}
			NIP.Mu.Unlock()
		}
	}

}

func notificationHandler(nr *NotificationRequest, sc <-chan bool) {

	var timerChan <-chan time.Time
	var tickerChan <-chan time.Time

	uuid := nr.Plan.ID.String()
	log.Println("[", uuid, "]", "Initiating notification plan")

	for n, s := range nr.Plan.Steps {

		u, err := url.Parse(s.Method)
		if err != nil {
			log.Println("Error parsing URI:", err)
			log.Println("Advancing to next step in plan.")
			continue
		}

		log.Println("[", uuid, "]", "Method:", s.Method)

		switch u.Scheme {
		case "phone":
			MakePhoneCall(u.Host, nr.Content, uuid)
		case "sms":
			SendSMS(u.Host, nr.Content, uuid, false)
		case "email":
			SendEmail(fmt.Sprint(u.User, "@", u.Host), nr.Content, uuid)
		}

		if n == len(nr.Plan.Steps)-1 {
			// Last step, so we use a Ticker and NotifyEveryPeriod
			tickerChan = time.NewTicker(s.NotifyEveryPeriod).C
			log.Println("[", uuid, "]", "Scheduling the next retry in", strconv.FormatFloat(s.NotifyEveryPeriod.Minutes(), 'f', 1, 64), "minutes")

		} else {
			// Not the last step, so we use a Timer and NotifyUntilPeriod
			timerChan = time.NewTimer(s.NotifyUntilPeriod).C
			log.Println("[", uuid, "]", "Scheduling the next notification step in", strconv.FormatFloat(s.NotifyUntilPeriod.Minutes(), 'f', 1, 64), "minutes")
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
				log.Println("[", uuid, "]", "Stop request received.  Terminating notifications.")
				NIP.Mu.Lock()
				defer NIP.Mu.Unlock()
				delete(NIP.Stoppers, uuid)
				delete(NIP.Messages, uuid)
				return
			}
		}

		log.Println("Loop broken")
	}
}
