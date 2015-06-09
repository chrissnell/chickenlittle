package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Username            string `yaml:"username" json:"username"`
	FullName            string `yaml:"full_name" json:"fullname"`
	VictorOpsRoutingKey string `yaml:"victorops_routing_key" json:"victorops_routing_key,omitempty"`
}

func (p *Person) Marshal() ([]byte, error) {
	jp, err := json.Marshal(&p)
	return jp, err
}

func (p *Person) Unmarshal(jp string) error {
	err := json.Unmarshal([]byte(jp), &p)
	return err
}

func (c *ChickenLittle) GetPerson(p string) (*Person, error) {
	jp, err := c.DB.Fetch("people", p)
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

func (c *ChickenLittle) GetAllPeople() ([]*Person, error) {
	var peeps []*Person

	jp, err := c.DB.FetchAll("people")
	if err != nil {
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

func (c *ChickenLittle) StorePerson(p *Person) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal person %+v", p)
	}

	err = c.DB.Store("people", p.Username, string(jp))
	if err != nil {
		return err
	}

	return nil
}

func (c *ChickenLittle) DeletePerson(p string) error {
	err := c.DB.Delete("people", p)
	if err != nil {
		return err
	}

	return nil
}
