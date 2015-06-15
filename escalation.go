package main

import (
	"encoding/json"
	"time"
)

type EscalationMethod int

const (
	NotifyOnDuty         EscalationMethod = iota // 0
	NotifyNextInRotation                         // 1
	NotifyOtherPerson                            // 2
	NotifyWebhook                                // 3
	NotifyEmail                                  // 4
)

// EscalationStep defines how alerts are escalated when an contact
// does not respond in time. Target takes a different meaning depending on
// the EscalationMethod. It is ignored on NotifyOnDuty or NotifyNextInRotation.
// For NotifyOtherPerson it's the name of another contact, for CallWebhook it's an URL
// and for SendEmail it's an email address.
type EscalationStep struct {
	TimeBeforeEscalation time.Duration    `yaml:"timebefore" json:"timebefore"` // How long to try the current step
	Method               EscalationMethod `yaml:"method" json:"method"`         // What action to take during this step
	Target               string           `yaml:"target" json:"target"`         // Who or what to do the action with
}

func (e *EscalationMethod) Marshal() ([]byte, error) {
	je, err := json.Marshal(&e)
	return je, err
}

func (e *EscalationMethod) Unmarshal(je string) error {
	err := json.Unmarshal([]byte(je), &e)
	return err
}

func (e *EscalationStep) Marshal() ([]byte, error) {
	je, err := json.Marshal(&e)
	return je, err
}

func (e *EscalationStep) Unmarshal(je string) error {
	err := json.Unmarshal([]byte(je), &e)
	return err
}
