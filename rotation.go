package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func (rp *RotationPolicy) Marshal() ([]byte, error) {
	jrp, err := json.Marshal(&rp)
	return jrp, err
}

func (rp *RotationPolicy) Unmarshal(jrp string) error {
	err := json.Unmarshal([]byte(jrp), &rp)
	return err
}

// Fetch a Rotation Policy from the DB
func (c *ChickenLittle) GetRotationPolicy(rp string) (*RotationPolicy, error) {
	jrp, err := c.DB.Fetch("rotationpolicies", rp)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch rotation policy %v from DB", rp)
	}

	policy := &RotationPolicy{}

	err = policy.Unmarshal(jrp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal rotation policy from DB.  Err: %v  JSON: %v", err, jrp)
	}

	return policy, nil
}

// Fetch every Rotation Policy from the DB
func (c *ChickenLittle) GetAllRotationPolicies() ([]*RotationPolicy, error) {
	var policies []*RotationPolicy

	jrp, err := c.DB.FetchAll("rotationpolicies")
	if err != nil {
		log.Println("Error fetching all rotation policies from DB:", err, "(Have you added any rotation policies?)")
		return nil, fmt.Errorf("Could not fetch all rotation policies from DB")
	}

	for _, v := range jrp {
		policy := &RotationPolicy{}

		err = policy.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal rotation policy from DB.  Err: %v  JSON: %v", err, jrp)
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// Store a Rotation Policy in the DB
func (c *ChickenLittle) StoreRotationPolicy(rp *RotationPolicy) error {
	jrp, err := rp.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal rotation policy %+v", rp)
	}

	// Note: the UUID needs to be generated at time of creation.  Do we generate it
	//       in this file or do we generate it as part of the rotation policy API?
	err = c.DB.Store("rotationpolicies", rp.UUID, string(jrp))
	if err != nil {
		return err
	}

	return nil
}

// Delete a Rotation Policy from the DB
func (c *ChickenLittle) DeleteRotationPolicy(rp string) error {
	err := c.DB.Delete("rotationpolicies", rp)
	if err != nil {
		return err
	}

	return nil
}
