package notification

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

func (e *Engine) notificationHandler(uuid string) {
	n, found := e.notifications[uuid]
	// TODO error checking
	_ = n
	_ = found
	sc, found := e.stoppers[uuid]
	_ = sc
	_ = found

	// iterate over every escalation step in the escalation plan until we run out of steps
	//for escStepNum, escStep := range n.Steps() {
	//	err := e.handleEscalationStep(n, escStep)
	//	if err != nil {
	//		log.Println("[" + uuid + "] Loop broken")
	//		break
	//	}
	//}
}

func (e *Engine) handleEscalationStep(n Notification, es EscalationStep) error {
	for notStepNum, notStep := range es.Steps() {
		// TODO repeat every period until we run out of time ...
		e.notify(notStep.Target(), n)
		_ = notStepNum // XXX
	}
	return nil
}

func (e *Engine) notify(target *url.URL, n Notification) error {
	switch target.Scheme {
	case "phone":
		e.MakePhoneCall(target.Host, n.Message(), n.ID())
	case "sms":
		e.SendSMS(target.Host, n.Message(), n.ID(), false)
	case "email":
		e.SendEmail(fmt.Sprintf("%s@%s", target.User, target.Host), n.Message(), n.ID())
	case "http":
		e.CallHTTP(target, n.Subject(), n.Message(), n.ID())
	case "https":
		e.CallHTTP(target, n.Subject(), n.Message(), n.ID())
	case "noop":
		log.Println("[" + n.ID() + "] Noop Notification: " + target.String())
	}
	// TODO error handling
	return nil
}

// notificationHandler receives notification requests from the engine and steps through the plan, making
// phone calls, sendins SMS, email, etc. as necessary
func (e *Engine) notificationPersonHandler(uuid string) {
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
