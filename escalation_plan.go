package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func (e *EscalationPlan) Marshal() ([]byte, error) {
	je, err := json.Marshal(&e)
	return je, err
}

func (e *EscalationPlan) Unmarshal(je string) error {
	err := json.Unmarshal([]byte(je), &e)
	return err
}

// Fetch an Escalation Plan from the DB
func (c *ChickenLittle) GetEscalationPlan(p string) (*EscalationPlan, error) {
	jp, err := c.DB.Fetch("escalationplans", p)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch escalation plan %v from DB", p)
	}

	plan := &EscalationPlan{}

	err = plan.Unmarshal(jp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal escalation plan from DB.  Err: %v  JSON: %v", err, jp)
	}

	return plan, nil
}

// Fetch every Escalation Plan from the DB
func (c *ChickenLittle) GetAllEscalationPlans() ([]*EscalationPlan, error) {
	var plans []*EscalationPlan

	jp, err := c.DB.FetchAll("escalationplans")
	if err != nil {
		log.Println("Error fetching all escalation plans from DB:", err, "(Have you added any escalation plans?)")
		return nil, fmt.Errorf("Could not fetch all escalation plans from DB")
	}

	for _, v := range jp {
		plan := &EscalationPlan{}

		err = plan.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal escalation plan from DB.  Err: %v  JSON: %v", err, jp)
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

// Store an Escalation Plan in the DB
func (c *ChickenLittle) StoreEscalationPlan(p *EscalationPlan) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal escalation plan %+v", p)
	}

	// Note: the UUID needs to be generated at time of creation.  Do we generate it
	//       in this file or do we generate it as part of the escalation plan API?
	err = c.DB.Store("escalationplans", p.UUID, string(jp))
	if err != nil {
		return err
	}

	return nil
}

// Delete an Escalation Plan from the DB
func (c *ChickenLittle) DeleteEscalationPlan(p string) error {
	err := c.DB.Delete("escalationplans", p)
	if err != nil {
		return err
	}

	return nil
}
