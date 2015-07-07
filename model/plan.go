package model

import (
	"net/url"
	"time"
)

// NotificationSteper is one step of a notification process
type NotificationSteper interface {
	Target() *url.URL
	Until() time.Duration
	Frequency() time.Duration
}

// EscalationSteper is on escalation step of a notification
type EscalationSteper interface {
	Steps() []NotificationSteper
}

// Notifier is a complete notification
type Notifier interface {
	Steps() []EscalationSteper
	ID() string
	Message() string
	Subject() string
	Stopper() chan struct{}
}
