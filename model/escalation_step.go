package model

import "time"

// EscalationMethod is one of the following escalation methods for use in an EscalationStep
type EscalationMethod int

// EscalationMethods as used in the EscalationSteps
const (
	NotifyOnDuty         EscalationMethod = iota // 0 - Notifies the on-duty person.  The first step of a plan typically uses this method
	NotifyNextInRotation                         // 1 - Notifies the person whose on-duty shift succeeds the current on-duty person.
	NotifyOtherPerson                            // 2 - Notifies another person, not necessarily part of the current rotation or team.  Could be a manager.
	NotifyWebhook                                // 3 - Calls a webhook (URL stored in EscalationStep.Target)
	NotifyEmail                                  // 4 - Notifies an e-mail address (e-mail address stored in EscalationStep.Target)
	NotifyAllInRotation                          // 5 - Notifies all persons in the current team rotation
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
