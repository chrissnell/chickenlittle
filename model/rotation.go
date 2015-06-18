package model

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

// Marshal implements the json Encoder interface
func (rp *RotationPolicy) Marshal() ([]byte, error) {
	jrp, err := json.Marshal(&rp)
	return jrp, err
}

// Unmarshal implements the json Decoder interface
func (rp *RotationPolicy) Unmarshal(jrp string) error {
	err := json.Unmarshal([]byte(jrp), &rp)
	return err
}

// GetRotationPolicy will fetch a Rotation Policy from the DB
func (m *Model) GetRotationPolicy(uuid string) (*RotationPolicy, error) {
	jrp, err := m.db.Fetch("rotationpolicies", uuid)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch rotation policy %v from DB", uuid)
	}

	policy := &RotationPolicy{}

	err = policy.Unmarshal(jrp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal rotation policy from DB.  Err: %v  JSON: %v", err, jrp)
	}

	return policy, nil
}

// GetAllRotationPolicies will fetch every Rotation Policy from the DB
func (m *Model) GetAllRotationPolicies() ([]*RotationPolicy, error) {
	var policies []*RotationPolicy

	jrp, err := m.db.FetchAll("rotationpolicies")
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

// StoreRotationPolicy will store a Rotation Policy in the DB
func (m *Model) StoreRotationPolicy(rp *RotationPolicy) error {
	jrp, err := rp.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal rotation policy %+v", rp)
	}

	// Note: the UUID needs to be generated at time of creation.  Do we generate it
	//       in this file or do we generate it as part of the rotation policy API?
	err = m.db.Store("rotationpolicies", rp.UUID, string(jrp))
	if err != nil {
		return err
	}

	return nil
}

// DeleteRotationPolicy will delete a Rotation Policy from the DB
func (m *Model) DeleteRotationPolicy(uuid string) error {
	err := m.db.Delete("rotationpolicies", uuid)
	if err != nil {
		return err
	}

	return nil
}
