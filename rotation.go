package main

import (
	"encoding/json"
	"time"
)

// RotationPolicy defines how shifts are rotated. If the frequency is zero
// no automatic rotations should be attempted.
type RotationPolicy struct {
	UUID              string        `yaml:"uuid" json:"uuid"`
	Description       string        `yaml:"description" json:"description"`
	RotationFrequency time.Duration `yaml:"frequency" json:"frequency"`
	RotateTime        time.Time     `yaml:"time" json:"time"`
}

func (r *RotationPolicy) Marshal() ([]byte, error) {
	jr, err := json.Marshal(&r)
	return jr, err
}

func (r *RotationPolicy) Unmarshal(jr string) error {
	err := json.Unmarshal([]byte(jr), &r)
	return err
}
