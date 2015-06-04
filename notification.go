package main

import (
	"encoding/json"
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
	Method            Method
	Data              string
	NotifyEveryPeriod time.Duration
	NotifyUntilPeriod time.Duration
}

type NotificationProdcedure struct {
	Username string
	Steps    []NotificationStep `json:",omitempty"`
}

func (np *NotificationProdcedure) Marshal() ([]byte, error) {
	jnp, err := json.Marshal(np)
	return jnp, err
}

func (np *NotificationProdcedure) Unmarshal(jnp string) error {
	err := json.Unmarshal([]byte(jnp), np)
	return err
}
