package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Username               string                 `yaml:"username"`
	FullName               string                 `yaml:"full_name"`
	VictorOpsRoutingKey    string                 `yaml:"victorops_routing_key" json:"victorops_routing_key,omitempty"`
	NotificationProdcedure NotificationProdcedure `json:"-"`
}

func (p *Person) Marshal() ([]byte, error) {
	jp, err := json.Marshal(&p)
	return jp, err
}

func (p *Person) Unmarshal(jp string) error {
	err := json.Unmarshal([]byte(jp), &p)
	return err
}

// Returns a santized Person struct without sensitive information
func (p *Person) Sanitized() *Person {
	var sp Person
	sp.FullName = p.FullName
	sp.Username = p.Username
	return &sp
}

func (c *ChickenLittle) RefreshPersonFromDB(p string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jp, err := c.DB.Fetch("people", p)
	if err != nil {
		return fmt.Errorf("Could not fetch person %v from DB", p)
	}

	peep := &Person{}

	err = peep.Unmarshal(jp)
	if err != nil {
		return fmt.Errorf("Could not unmarshal person from DB.  Err: %v  JSON: %v", err, jp)
	}

	c.People[p] = peep

	return nil
}

func (c *ChickenLittle) GetPerson(p string) (*Person, error) {

	err := c.RefreshPersonFromDB(p)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	peep, pres := c.People[p]

	if !pres {
		return nil, fmt.Errorf("No such person: %v", p)
	}

	return peep, nil
}

func (c *ChickenLittle) StorePerson(p *Person) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.People[p.Username] = p

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
