package main

import (
	"encoding/json"
	"time"
)

// EscalationPlan defines how team alerts are escalated when a contact does
// not respond in time.  It consists of a series of EscalationSteps that are
// taken, in series, until a team alert is acknowledged
type EscalationPlan struct {
	UUID        string           `json:"uuid"`
	Description string           `json:"description"`
	Steps       []EscalationStep `json:"escalation_steps,omitempty"`
}

type EscalationMethod int

const (
	NotifyOnDuty         EscalationMethod = iota // 0 - Notifies the on-duty person.  The first step of a plan typically uses this method
	NotifyNextInRotation                         // 1 - Notifies the person whose on-duty shift succeeds the current on-duty person.
	NotifyOtherPerson                            // 2 - Notifies another person, not necessarily part of the current rotation or team.  Could be a manager.
	NotifyWebhook                                // 3 - Calls a webhook (URL stored in EscalationStep.Target)
	NotifyEmail                                  // 4 - Notifies an e-mail address (e-mail address stored in EscalationStep.Target)
)

// EscalationStep defines how a particular stage of team notification is handled.
//
// TimeBeforeEscalation is how long the current EscalationStep waits for an acknowledgement
// before proceeding to the next EscalationStep.  Note: this time period is independent of
// any person's own notification plan.  Individual notification plans will still be followed
// but the EscalationSteps will continue if acknowledgement isn't recieved.
//
// Method is an EscalationMethod as defined in the iota above.
//
// Target takes a different meaning depending on the EscalationMethod.  See definitions above.
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
