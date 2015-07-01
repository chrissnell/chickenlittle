package model

import (
	"net/url"
	"time"
)

// notificationStep is a single notification step implementing NotificationSteper
type notificationStep struct {
	target      *url.URL
	repeatAfter time.Duration
	repeatUntil time.Duration
}

// escalationStep is a single escalation step implementing EscalationSteper
type escalationStep struct {
	steps []NotificationSteper
}

// notificationJob is a Notification implementing Notifier
type notificationJob struct {
	id      string
	steps   []EscalationSteper
	subject string
	message string
	stopper chan struct{}
}

// Target returns the notification target as valid URL
func (n *notificationStep) Target() *url.URL {
	return n.target
}

// Until returns the duration how long this step should be repeated
func (n *notificationStep) Until() time.Duration {
	return n.repeatUntil
}

// Frequency returns the duration after that the current step
// will be retried until the Until duration expires
// and the next step is attempted.
func (n *notificationStep) Frequency() time.Duration {
	return n.repeatAfter
}

// Steps returns all notification steps of this escalation step
func (e *escalationStep) Steps() []NotificationSteper {
	return e.steps
}

// Steps returns all escalation steps of this notification
func (n *notificationJob) Steps() []EscalationSteper {
	return n.steps
}

// ID will return the UUID of this notification
func (n *notificationJob) ID() string {
	return n.id
}

// Message will return the message (body) of this notification
func (n *notificationJob) Message() string {
	return n.message
}

// Subject will return the summary (subject) of this notification
func (n *notificationJob) Subject() string {
	return n.subject
}

// Stopper will return the stop channel of this notification
func (n *notificationJob) Stopper() chan struct{} {
	return n.stopper
}
