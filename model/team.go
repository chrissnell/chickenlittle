package model

import (
	"encoding/json"
	"fmt"
	"log"
)

// Team contains a team of people with a defined rotations policy and an escalation plan
type Team struct {
	Name           string   `yaml:"name" json:"name"`                       // The name of the team
	Description    string   `yaml:"description" json:"description"`         // a human readable description of this team
	Members        []string `yaml:"members" json:"members"`                 // a list of members, SHOULD be valid users from the people bucket
	RotationPolicy string   `yaml:"rotation_policy" json:"rotation_policy"` // the policy for automatically changing the notification order
	EscalationPlan string   `yaml:"escalation_plan" json:"escalation_plan"` // the current escalation plan for getting hold of a team member
}

// Marshal implements the json Encoder interface
func (t *Team) Marshal() ([]byte, error) {
	jt, err := json.Marshal(&t)
	return jt, err
}

// Unmarshal implements the json Decoder interface
func (t *Team) Unmarshal(jt string) error {
	err := json.Unmarshal([]byte(jt), &t)
	return err
}

// GetTeam will fetch a Team from the DB
func (m *Model) GetTeam(t string) (*Team, error) {
	jt, err := m.db.Fetch("teams", t)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch team %v from DB", t)
	}

	team := &Team{}

	err = team.Unmarshal(jt)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal team from DB.  Err: %v  JSON: %v", err, jt)
	}

	return team, nil
}

// GetAllTeams will fetch every Team from the DB
func (m *Model) GetAllTeams() ([]*Team, error) {
	var teams []*Team

	jt, err := m.db.FetchAll("teams")
	if err != nil {
		log.Println("Error fetching all teams from DB:", err, "(Have you added any teams?)")
		return nil, fmt.Errorf("Could not fetch all teams from DB")
	}

	for _, v := range jt {
		team := &Team{}

		err = team.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal team from DB.  Err: %v  JSON: %v", err, jt)
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// StoreTeam will store a Team in the DB
func (m *Model) StoreTeam(t *Team) error {
	jt, err := t.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal team %+v", t)
	}

	err = m.db.Store("teams", t.Name, string(jt))
	if err != nil {
		return err
	}

	return nil
}

// DeleteTeam will delete a Team from the DB
func (m *Model) DeleteTeam(t string) error {
	err := m.db.Delete("teams", t)
	if err != nil {
		return err
	}

	return nil
}
