package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/twinj/uuid"
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

// getNotificationSteps will return a fully realized (all indirection resolved) notification plan
// for a single user.
func (m *Model) getNotificationSteps(username string) ([]NotificationSteper, error) {
	steps := make([]NotificationSteper, 0, 1)
	// fetch notification plan from DB
	np, err := m.GetNotificationPlan(username)
	if err != nil {
		return steps, err
	}
	for _, step := range np.Steps {
		url, err := url.Parse(step.Method)
		if err != nil {
			log.Println("Skipping invalid notification step for user ", username)
			continue
		}
		ns := &notificationStep{
			target:      url,
			repeatAfter: step.NotifyEveryPeriod,
			repeatUntil: step.NotifyUntilPeriod,
		}
		steps = append(steps, ns)
	}
	return steps, nil
}

// GetNotificationForPerson creates a new notificationJob given a person (username), subject and message.
// A Notification for one person will contain an escalation plan containing exactly one step which
// consists of the notification steps from the notification plan of this user.
func (m *Model) GetNotificationForPerson(username string, subject string, message string) (notificationJob, error) {
	notSteps, err := m.getNotificationSteps(username)
	if err != nil {
		return notificationJob{}, err
	}

	// wrap the notification plan in a dummy escalation plan
	escStep := &escalationStep{
		steps: notSteps,
	}
	escSteps := make([]EscalationSteper, 0, 1)
	escSteps = append(escSteps, escStep)
	// Assign a UUID to this notification. The UUID is used to track notifications-in-progress and to stop
	// them when requested.
	uuid.SwitchFormat(uuid.CleanHyphen)
	n := notificationJob{
		id:      uuid.NewV4().String(),
		subject: subject,
		message: message,
		steps:   escSteps,
		stopper: make(chan struct{}),
	}

	return n, nil
}
