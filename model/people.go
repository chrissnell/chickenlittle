package model

import (
	"encoding/json"
	"fmt"
	"log"
)

// Person holds a single person. A Person can be part of a team or being notified directly, if it has an notification plan.
type Person struct {
	Username            string `yaml:"username" json:"username"`
	FullName            string `yaml:"full_name" json:"fullname"`
	VictorOpsRoutingKey string `yaml:"victorops_routing_key" json:"victorops_routing_key,omitempty"` // TODO(dschulz) perhaps this should be something more flexible, i.e. an map for storing arbitary keys
}

// Marshal implements the json Encoder interface
func (p *Person) Marshal() ([]byte, error) {
	jp, err := json.Marshal(&p)
	return jp, err
}

// Unmarshal implements the json Decoder interface
func (p *Person) Unmarshal(jp string) error {
	err := json.Unmarshal([]byte(jp), &p)
	return err
}

// GetPerson will fetch a Person from the DB
func (m *Model) GetPerson(p string) (*Person, error) {
	jp, err := m.db.Fetch("people", p)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch person %v from DB", p)
	}

	peep := &Person{}

	err = peep.Unmarshal(jp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal person from DB.  Err: %v  JSON: %v", err, jp)
	}

	return peep, nil
}

// GetAllPeople will fetch every Person from the DB
func (m *Model) GetAllPeople() ([]*Person, error) {
	var peeps []*Person

	jp, err := m.db.FetchAll("people")
	if err != nil {
		log.Println("Error fetching all people from DB:", err, "(Have you added any people?)")
		return nil, fmt.Errorf("Could not fetch all people from DB")
	}

	for _, v := range jp {
		peep := &Person{}

		err = peep.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal person from DB.  Err: %v  JSON: %v", err, jp)
		}

		peeps = append(peeps, peep)
	}

	return peeps, nil
}

// StorePerson will store a Person in the DB
func (m *Model) StorePerson(p *Person) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal person %+v", p)
	}

	err = m.db.Store("people", p.Username, string(jp))
	if err != nil {
		return err
	}

	return nil
}

// DeletePerson will delete a Person from the DB
func (m *Model) DeletePerson(p string) error {
	err := m.db.Delete("people", p)
	if err != nil {
		return err
	}

	return nil
}
