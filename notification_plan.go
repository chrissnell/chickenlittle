package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	Method uint8
)

const (
	Voice Method = iota
	SMS
	Email
	Callback
)

type NotificationStep struct {
	Method            Method        `json:"method"`
	Data              string        `json:"data"`
	NotifyEveryPeriod time.Duration `json:"notify_every_period"`
	NotifyUntilPeriod time.Duration `json:"notify_until_period"`
}

type NotificationPlan struct {
	Username string             `json:"username"`
	Steps    []NotificationStep `json:"steps,omitempty"`
}

func (np *NotificationPlan) Marshal() ([]byte, error) {
	jnp, err := json.Marshal(np)
	return jnp, err
}

func (np *NotificationPlan) Unmarshal(jnp string) error {
	err := json.Unmarshal([]byte(jnp), np)
	return err
}

func (c *ChickenLittle) GetNotificationPlan(username string) (*NotificationPlan, error) {
	jp, err := c.DB.Fetch("notificationplans", username)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch notification plan from DB: plan for %v does not exist", username)
	}

	plan := &NotificationPlan{}

	err = plan.Unmarshal(jp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal notification plan from DB.  Err: %v  JSON: %v", err, jp)
	}

	return plan, nil
}

func (c *ChickenLittle) StoreNotificationPlan(p *NotificationPlan) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal person %+v", p)
	}

	err = c.DB.Store("notificationplans", p.Username, string(jp))
	if err != nil {
		return err
	}

	return nil
}

func (c *ChickenLittle) DeleteNotificationPlan(username string) error {
	err := c.DB.Delete("notificationplans", username)
	if err != nil {
		return err
	}

	return nil
}
