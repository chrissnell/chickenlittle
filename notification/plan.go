package notification

import (
	"net/url"
	"time"
)

// NotificationStep is one step of a notification process
type NotificationStep interface {
	NotifyMethod() string
	Frequency() time.Duration
	Until() time.Duration
	Target() *url.URL
}

type EscalationStep interface {
	Steps() []NotificationStep
	ID() string
}

type Notification interface {
	NextStep() (NotificationStep, bool)
	ID() string
	Message() string
	Subject() string
	Stopper() chan<- struct{}
}

// NotificationRequest is an interface consumed by the notification handler. It
// abstracts the differences between people and teams away.
type TeamNotification interface {
	Steps() []EscalationStep // returns the next notification step, handling possible escalations, or true if there are not more steps available
	ID() string              // return the assigned UUID
	Message() string         // return the message (content) to be sent
	Subject() string         // return a mesage subject, if availabe, of a configurable default
}

type notificationStep struct {
	target        *url.URL
	NextStepAfter time.Duration
	RepeatUntil   time.Duration
}
type escalationStep struct {
	notificationSteps []notificationStep
}
type escalationPlan struct {
	escalationSteps []escalationStep
}
type notification struct {
	id      string
	plan    escalationPlan
	subject string
	message string
	stopper chan struct{}
}

func (n *notificationStep) Target() *url.URL {
	return n.target
}

func NotificationForPerson(username, subject, message string) notification {
	notSteps := make([]notificationStep, 1)
	// TODO fetch notification plan from DB
	escStep := escalationStep{
		notificationSteps: notSteps,
	}
	escSteps := make([]escalationStep, 1)
	escSteps = append(escSteps, escStep)
	p := escalationPlan{
		escalationSteps: escSteps,
	}
	n := notification{
		id:      "uuid",
		subject: subject,
		message: message,
		plan:    p,
		stopper: make(chan struct{}),
	}

	return n
}

func NotificationForTeam(name, subject, message string) notification {
	// TODO fetch escalation plan and team from db
	escSteps := make([]escalationStep, 1)
	// TODO iterate over esc plan and create a step for every step
	p := escalationPlan{
		escalationSteps: escSteps,
	}
	n := notification{
		id:      "uuid",
		subject: subject,
		message: message,
		plan:    p,
		stopper: make(chan struct{}),
	}
	return n
}
