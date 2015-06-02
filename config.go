package main

import (
	"fmt"
)

type Config struct {
	People       []Person      `yaml:"people"`
	Service      ServiceConfig `yaml:"service"`
	Integrations Integrations  `yaml:"integrations"`
}

type Person struct {
	User                string `yaml:"user"`
	FullName            string `yaml:"full_name"`
	VictorOpsRoutingKey string `yaml:"victorops_routing_key" json:"victorops_routing_key,omitempty"`
}

type ServiceConfig struct {
	ListenAddr string `yaml:"listen_address"`
}

type Integrations struct {
	HipChat   HipChat   `yaml:"hipchat"`
	VictorOps VictorOps `yaml:"victorops"`
}

type VictorOps struct {
	APIKey string `yaml:"api_key"`
}

type HipChat struct {
	HipChatAuthToken    string `yaml:"hipchat_auth_token"`
	HipChatAnnounceRoom string `yaml:"hipchat_announce_room"`
}

// Returns a santized Person struct without sensitive information
func (p *Person) Sanitized() *Person {
	var sp Person
	sp.FullName = p.FullName
	sp.User = p.User
	return &sp
}

func (c *Config) PeopleAsMap() map[string]*Person {
	pmap := make(map[string]*Person)
	for _, v := range c.People {
		pmap[v.User] = v.Sanitized()
	}

	return pmap
}

func (c *Config) GetPerson(p string) (*Person, error) {
	for _, v := range c.People {
		if v.User == p {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("No such person: %v", p)
}
