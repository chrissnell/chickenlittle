package model

import (
	"encoding/json"
	"fmt"
	"log"
)

// EscalationPlan defines how team alerts are escalated when a contact does
// not respond in time.  It consists of a series of EscalationSteps that are
// taken, in series, until a team alert is acknowledged
type EscalationPlan struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Steps       []EscalationStep `json:"escalation_steps,omitempty"`
}

// Marshal implements the json Encoder interface
func (e *EscalationPlan) Marshal() ([]byte, error) {
	je, err := json.Marshal(&e)
	return je, err
}

// Unmarshal implements the json Decoder interface
func (e *EscalationPlan) Unmarshal(je string) error {
	err := json.Unmarshal([]byte(je), &e)
	return err
}

// GetEscalationPlan will fetch an Escalation Plan from the DB
func (m *Model) GetEscalationPlan(p string) (*EscalationPlan, error) {
	jp, err := m.db.Fetch("escalationplans", p)
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

// GetAllEscalationPlans will fetch every Escalation Plan from the DB
func (m *Model) GetAllEscalationPlans() ([]*EscalationPlan, error) {
	var plans []*EscalationPlan

	jp, err := m.db.FetchAll("escalationplans")
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

// StoreEscalationPlan will store an Escalation Plan in the DB
func (m *Model) StoreEscalationPlan(p *EscalationPlan) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal escalation plan %+v", p)
	}

	// Note: the UUID needs to be generated at time of creation.  Do we generate it
	//       in this file or do we generate it as part of the escalation plan API?
	err = m.db.Store("escalationplans", p.Name, string(jp))
	if err != nil {
		return err
	}

	return nil
}

// DeleteEscalationPlan will delete an Escalation Plan from the DB
func (m *Model) DeleteEscalationPlan(p string) error {
	err := m.db.Delete("escalationplans", p)
	if err != nil {
		return err
	}

	return nil
}
