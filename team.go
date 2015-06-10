package main

import (
	"encoding/json"
	"fmt"
)

type Team struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Members     []string `yaml:"members" json:"members"`
}

func (t *Team) Marshal() ([]byte, error) {
	jt, err := json.Marshal(&t)
	return jt, err
}

func (t *Team) Unmarshal(jt string) error {
	err := json.Unmarshal([]byte(jt), &t)
	return err
}

// Fetch a Team from the DB
func (c *ChickenLittle) GetTeam(t string) (*Team, error) {
	jt, err := c.DB.Fetch("teams", t)
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

// Fetch every Team from the DB
func (c *ChickenLittle) GetAllTeams() ([]*Team, error) {
	var teams []*Team

	jt, err := c.DB.FetchAll("teams")
	if err != nil {
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

// Store a Team in the DB
func (c *ChickenLittle) StoreTeam(t *Team) error {
	jt, err := t.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal team %+v", t)
	}

	err = c.DB.Store("teams", t.Name, string(jt))
	if err != nil {
		return err
	}

	return nil
}

// Delete a Team from the DB
func (c *ChickenLittle) DeleteTeam(t string) error {
	err := c.DB.Delete("teams", t)
	if err != nil {
		return err
	}

	return nil
}
