package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// NotificationPlan is a list of steps to take to notify on person.
type NotificationPlan struct {
	Username string             `json:"username"`        // the username, must be a valid user from the bucket "people".
	Steps    []NotificationStep `json:"steps,omitempty"` // the steps to take
}

// NotificationStep is a single step to take during the notification for one user.
type NotificationStep struct {
	Method            string        `json:"method"`              // the method to perform, must be a valid URL like sms://5551234
	NotifyEveryPeriod time.Duration `json:"notify_every_period"` // how often to repeat this notification step
	NotifyUntilPeriod time.Duration `json:"notify_until_period"` // how long to try to notify during this step
}

// NotifyMethod will return the notification method
func (ns NotificationStep) NotifyMethod() string {
	return ns.Method
}

// Frequency will return the notification frequency
func (ns NotificationStep) Frequency() time.Duration {
	return ns.NotifyEveryPeriod
}

// Until will return the length of the notification period
func (ns NotificationStep) Until() time.Duration {
	return ns.NotifyUntilPeriod
}

// Marshal implements the json Encoder interface
func (np *NotificationPlan) Marshal() ([]byte, error) {
	jnp, err := json.Marshal(np)
	return jnp, err
}

// Unmarshal implements the json Decoder interface
func (np *NotificationPlan) Unmarshal(jnp string) error {
	err := json.Unmarshal([]byte(jnp), np)
	return err
}

// GetNotificationPlan will fetch a NotificationPlan from the DB
func (m *Model) GetNotificationPlan(username string) (*NotificationPlan, error) {
	jp, err := m.db.Fetch("notificationplans", username)
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

// StoreNotificationPlan will store a NotificationPlan in the DB
func (m *Model) StoreNotificationPlan(p *NotificationPlan) error {
	jp, err := p.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal person %+v", p)
	}

	err = m.db.Store("notificationplans", p.Username, string(jp))
	if err != nil {
		return err
	}

	return nil
}

// DeleteNotificationPlan will delete a NotificationPlan from the DB
func (m *Model) DeleteNotificationPlan(username string) error {
	err := m.db.Delete("notificationplans", username)
	if err != nil {
		return err
	}

	return nil
}
