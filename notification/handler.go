package notification

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/chrissnell/chickenlittle/model"
)

// notificationHandler receives notification requests from the engine and steps through the plan, making
// phone calls, sendins SMS, email, etc. as necessary
func (e *Engine) notificationHandler(uuid string) {
	// look up notification
	e.mutex.Lock()
	n, found := e.notifications[uuid]
	e.mutex.Unlock()
	if !found || n == nil {
		log.Println("[" + uuid + "] No notification found")
		return
	}

	// iterate over every escalation step in the escalation plan until we run out of steps
	for _, escStep := range n.Steps() {
		err := e.handleEscalationStep(n, escStep)
		if err != nil {
			log.Println("[", uuid, "] Loop broken")
			break
		}
	}
	log.Println("[", n.ID(), "]", "Exhausted all escalation steps. Removing notification.")
	e.mutex.Lock()
	delete(e.notifications, uuid)
	e.mutex.Unlock()
}

func (e *Engine) handleEscalationStep(n model.Notifier, es model.EscalationSteper) error {
STEPLOOP:
	for _, notStep := range es.Steps() {
		stepTimeout := time.After(notStep.Until())
	TIMERLOOP:
		for {
			retryTimeout := time.After(notStep.Frequency())
			err := e.notify(notStep.Target(), n)
			if err != nil {
				log.Printf("[", n.ID(), "] Notification failed: %s", err)
			}
			select {
			case <-stepTimeout:
				// Our timer for this step has expired so we break the outer loop to proceed to the next step.
				log.Println("[", n.ID(), "]", "Step timer expired.  Proceeding to next plan step.")
				continue STEPLOOP
			case <-retryTimeout:
				// Our ticker for this step expired, so we'll break the inner loop and try this step again.
				log.Println("[", n.ID(), "]", "**Tick**  Retry contact method!")
				log.Println("[", n.ID(), "]", "Waiting", strconv.FormatFloat(notStep.Frequency().Seconds(), 'f', 1, 64), "seconds")
				continue TIMERLOOP
			case <-n.Stopper():
				log.Println("[", n.ID(), "]", "Stop request received.  Terminating notifications.")
				e.mutex.Lock()
				defer e.mutex.Unlock()
				delete(e.notifications, n.ID())
				return errors.New("Stop request received")
			}
		}
	}
	log.Println("[", n.ID(), "]", "Exhausted all notification steps. Escalating to next step.")
	return nil
}

func (e *Engine) notify(target *url.URL, n model.Notifier) error {
	var err error

	log.Println("[ " + n.ID() + " ] Sending notification to " + target.String())

	switch target.Scheme {
	case "phone":
		err = e.MakePhoneCall(target.Host, n.Message(), n.ID())
	case "sms":
		err = e.SendSMS(target.Host, n.Message(), n.ID(), false)
	case "email":
		err = e.SendEmail(fmt.Sprintf("%s@%s", target.User, target.Host), n.Message(), n.ID())
	case "http":
		err = e.CallHTTP(target, n.Subject(), n.Message(), n.ID())
	case "https":
		err = e.CallHTTP(target, n.Subject(), n.Message(), n.ID())
	case "noop":
		log.Println("[ " + n.ID() + " ] Noop Notification: " + target.String())
	default:
		log.Println("[ " + n.ID() + " ] Unknown Notification scheme: " + target.String())
	}

	return err
}
