package notification

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/chrissnell/chickenlittle/config"
)

// NotificationStep is one step of a notification process
type NotificationStep interface {
	NotifyMethod() string
	Frequency() time.Duration
	Until() time.Duration
}

// NotificationRequest is an interface consumed by the notification handler. It
// abstracts the differences between people and teams away.
type Notification interface {
	NextStep() (NotificationStep, bool) // returns the next notification step, handling possible escalations, or true if there are not more steps available
	ID() string                         // return the assigned UUID
	Message() string                    // return the message (content) to be sent
	Subject() string                    // return a mesage subject, if availabe, of a configurable default
}

// Engine is the core of the notification engine handling all notifications
// and feedback handling.
type Engine struct {
	Config        config.Config
	planChan      chan Notification
	stopChan      chan string
	mutex         *sync.Mutex // protects all below
	stoppers      map[string]chan struct{}
	notifications map[string]Notification
	conversations map[string]string
}

// New creates a new notification engine instace with a running worker goroutine.
func New(c config.Config) *Engine {
	ne := &Engine{
		Config:        c,
		planChan:      make(chan Notification, 100),
		stopChan:      make(chan string, 100),
		mutex:         &sync.Mutex{},
		stoppers:      make(map[string]chan struct{}),
		notifications: make(map[string]Notification),
		conversations: make(map[string]string),
	}
	go ne.start()
	return ne
}

// start should be run exactly once from within the constructor.
// It will listen on the embedded chans for incoming notifications
// and stop requests and launch the appropriate methods.
func (e *Engine) start() {
	for {
		select {
		case notifyReq := <-e.planChan:
			e.startNotification(notifyReq)
		case uuid := <-e.stopChan:
			e.stopNotification(uuid)
		}
	}
}

// GetMessage will return the message for the notification
// identified by the given UUID or an empty string.
func (e *Engine) GetMessage(uuid string) string {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if nr, found := e.notifications[uuid]; found {
		return nr.Message()
	}
	return ""
}

// SetConversation will set the value to the given key in the
// conversation map. Value should be an UUID string.
func (e *Engine) SetConversation(key, value string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.conversations[key] = value
}

// GetConversation will return the UUID for the given
// conversation key or an error.
func (e *Engine) GetConversation(key string) (string, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if value, found := e.conversations[key]; found {
		return value, nil
	}
	return "", errors.New("Conversation not found")
}

// RemoveConversation will remove the given conversation
// for the conversation map.
func (e *Engine) RemoveConversation(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.conversations, key)
}

// EnqueueNotification will place a new notification in the
// notification queue chan. It will eventually be picked up by the embedded
// worker goroutine and processed in it's own handler.
func (e *Engine) EnqueueNotification(nr Notification) {
	e.planChan <- nr
}

// CancelNotification will place a stop request in the stop queue chan.
// It will be picked up by a running notification handler and the handler
// should quit shortly after.
func (e *Engine) CancelNotification(uuid string) {
	e.stopChan <- uuid
}

// startNotification will setup a new notification handler for the
// given notification request and run it in it's own goroutine.
func (e *Engine) startNotification(nr Notification) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	id := nr.ID()
	// Create a new stopper channel for this plan
	e.stoppers[id] = make(chan struct{})
	// Store the message in the message map
	e.notifications[id] = nr

	// asyncronously launch the notification plan processing
	go e.notificationHandler(id)
}

// IsNotification looks up the given notification and
// returns true if it is a notification currently being executed.
func (e *Engine) IsNotification(uuid string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, found := e.stoppers[uuid]
	return found
}

// stopNotification will send a stop request to the given notification
// or silently fail if the uuid is not a running notification.
func (e *Engine) stopNotification(uuid string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// check if there are notifications going for this UUID
	stopper, found := e.stoppers[uuid]
	if !found {
		return
	}
	log.Println("[", uuid, "]", "Sending a stop notification to the plan processor")

	// It's in progress, so we'll send a message on its Stopper to
	// be received by the goroutine executing the plan
	stopper <- struct{}{}
}

// notificationHandler receives notification requests from the engine and steps through the plan, making
// phone calls, sendins SMS, email, etc. as necessary
func (e *Engine) notificationHandler(uuid string) {
	var timerChan <-chan time.Time
	var tickerChan <-chan time.Time

	// look up notification
	e.mutex.Lock()
	n, found := e.notifications[uuid]
	e.mutex.Unlock()
	if !found {
		return
	}
	// look up stop chan
	e.mutex.Lock()
	sc, found := e.stoppers[uuid]
	e.mutex.Unlock()
	if !found {
		return
	}

	log.Println("[", uuid, "]", "Initiating notification plan")

	// Iterate through each step of the plan
	for {
		s, lastStep := n.NextStep()

		u, err := url.Parse(s.NotifyMethod())
		if err != nil {
			log.Println("Error parsing URI:", err)
			log.Println("Advancing to next step in plan.")
			continue
		}
		log.Println("[", uuid, "]", "Method:", s.NotifyMethod())

		// This outer loop repeats a notification until it's acknowledged.  It can be broken by the expiration of the timer for this step,
		// or by a stop request.
	stepLoop:
		for {
			// Take the appropriate action, depending on the type of notification
			switch u.Scheme {
			case "phone":
				e.MakePhoneCall(u.Host, n.Message(), uuid)
			case "sms":
				e.SendSMS(u.Host, n.Message(), uuid, false)
			case "email":
				e.SendEmail(fmt.Sprint(u.User, "@", u.Host), n.Message(), uuid)
			}
			if lastStep {
				// We're at the last step of the plan, so this step will repeat until ackknowledged. We use a Ticker and set its period to NotifyEveryPeriod
				tickerChan = time.NewTicker(s.Frequency()).C
				log.Println("[", uuid, "]", "Scheduling the next retry in", strconv.FormatFloat(s.Frequency().Minutes(), 'f', 1, 64), "minutes")

			} else {
				// We're not at the last step, so we only run this step once.  We use Timer set its duration to NotifyUntilPeriod
				timerChan = time.NewTimer(s.Until()).C
				log.Println("[", uuid, "]", "Scheduling the next notification step in", strconv.FormatFloat(s.Until().Minutes(), 'f', 1, 64), "minutes")
			}
			// This inner loop selects over various channels to receive timers and stop requests.
			// It can be broken by an expiring Timer (signaling that it's time to proceed to the next step) or a stop request.
		timerLoop:
			for {
				select {
				case <-timerChan:
					// Our timer for this step has expired so we break the outer loop to proceed to the next step.
					log.Println("[", uuid, "]", "Step timer expired.  Proceeding to next plan step.")
					break stepLoop
				case <-tickerChan:
					// Our ticker for this step expired, so we'll break the inner loop and try this step again.
					log.Println("[", uuid, "]", "**Tick**  Retry contact method!")
					log.Println("[", uuid, "]", "Waiting", strconv.FormatFloat(s.Frequency().Minutes(), 'f', 1, 64), "minutes]")
					break timerLoop
				case <-sc:
					log.Println("[", uuid, "]", "Stop request received.  Terminating notifications.")
					e.mutex.Lock()
					defer e.mutex.Unlock()
					delete(e.stoppers, uuid)
					delete(e.notifications, uuid)
					return
				}
			}

			log.Println("Loop broken")
		}
	}
}
