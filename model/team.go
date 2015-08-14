package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/twinj/uuid"
)

// Team contains a team of people with a defined rotations policy and an escalation plan
type Team struct {
	Name           string   `yaml:"name" json:"name"`                       // The name of the team
	Description    string   `yaml:"description" json:"description"`         // a human readable description of this team
	Members        []string `yaml:"members" json:"members"`                 // a list of members, SHOULD be valid users from the people bucket
	RotationPolicy string   `yaml:"rotation_policy" json:"rotation_policy"` // the policy for automatically changing the notification order
	EscalationPlan string   `yaml:"escalation_plan" json:"escalation_plan"` // the current escalation plan for getting hold of a team member
}

// Marshal implements the json Encoder interface
func (t *Team) Marshal() ([]byte, error) {
	jt, err := json.Marshal(&t)
	return jt, err
}

// Unmarshal implements the json Decoder interface
func (t *Team) Unmarshal(jt string) error {
	err := json.Unmarshal([]byte(jt), &t)
	return err
}

// GetTeam will fetch a Team from the DB
func (m *Model) GetTeam(t string) (*Team, error) {
	jt, err := m.db.Fetch("teams", t)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch team %v from DB", t)
	}

	team := &Team{}

	err = team.Unmarshal(jt)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal team from DB.  Err: %v  JSON: %v", err, jt)
	}

	return team, nil
}

// GetAllTeams will fetch every Team from the DB
func (m *Model) GetAllTeams() ([]*Team, error) {
	var teams []*Team

	jt, err := m.db.FetchAll("teams")
	if err != nil {
		log.Println("Error fetching all teams from DB:", err, "(Have you added any teams?)")
		return nil, fmt.Errorf("Could not fetch all teams from DB")
	}

	for _, v := range jt {
		team := &Team{}

		err = team.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal team from DB.  Err: %v  JSON: %v", err, jt)
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// StoreTeam will store a Team in the DB
func (m *Model) StoreTeam(t *Team) error {
	jt, err := t.Marshal()
	if err != nil {
		return fmt.Errorf("Could not marshal team %+v", t)
	}

	err = m.db.Store("teams", t.Name, string(jt))
	if err != nil {
		return err
	}

	return nil
}

// DeleteTeam will delete a Team from the DB
func (m *Model) DeleteTeam(t string) error {
	err := m.db.Delete("teams", t)
	if err != nil {
		return err
	}

	return nil
}

// getEscalationSteps will return a fully realized (all indirections resolved) notification plan
// for a team.
func (m *Model) getEscalationSteps(name string) ([]EscalationSteper, error) {
	escSteps := make([]EscalationSteper, 0, 1)
	// fetch escalation plan and team from db
	escPlan, err := m.GetEscalationPlan(name)
	if err != nil {
		return escSteps, err
	}
	team, err := m.GetTeam(name)
	if err != nil {
		return escSteps, err
	}
	curMember := 0
	// iterate over esc plan and create a step for every step
STEPS:
	for n, step := range escPlan.Steps {
		es := &escalationStep{}
		switch step.Method {
		// NotifyOnDuty will always notify the first member in the team array. On shift rotation
		// this array will be reordered and the next one on duty will be put on position 0.
		case NotifyOnDuty:
			ns, err := m.getNotificationSteps(team.Members[0])
			if err != nil {
				log.Printf("Error constructing escalation plan for %s in step %d: %s", name, n, err)
				continue STEPS
			}
			es.steps = ns
		// NotifyNextInRotation will start at the first team member after the one on duty and continue
		// notifying members at each occurency of NotifyNextInRotation until all team members are exhausted.
		// At that point further NotifyNextInRotatation steps will silently be ignored.
		case NotifyNextInRotation:
			curMember++
			if curMember >= len(team.Members) {
				log.Printf("No more team members available while constructing escalation plan for %s in step %d", name, n)
				continue STEPS
			}
			ns, err := m.getNotificationSteps(team.Members[curMember])
			if err != nil {
				log.Printf("Error constructing escalation plan for %s in step %d: %s", name, n, err)
				continue STEPS
			}
			es.steps = ns
		// NotifyAllInRotation will add all available team members to the escalation plan.
		case NotifyAllInRotation:
			if len(team.Members) < 2 {
				log.Printf("Not enough team mebers available whil constructing escalation plan for %s in step %d", name, n)
				continue STEPS
			}
			for i := 1; i < len(team.Members)-1; i++ {
				ns, err := m.getNotificationSteps(team.Members[i])
				if err != nil {
					continue
				}
				es := &escalationStep{}
				es.steps = ns
				escSteps = append(escSteps, es)
			}
			ns, err := m.getNotificationSteps(team.Members[len(team.Members)-1])
			if err != nil {
				continue STEPS
			}
			es.steps = ns
		// NotifyOtherPerson will create a notification step for a given person. This person must exist
		// somewhere in CL. In the current team or any other doesn't matter.
		case NotifyOtherPerson:
			ns, err := m.getNotificationSteps(step.Target)
			if err != nil {
				log.Printf("Error constructing escalation plan for %s in step %d: %s", name, n, err)
				continue STEPS
			}
			es.steps = ns
		// NotifyWebhook will notify make an HTTP POST request to the given URL.
		case NotifyWebhook:
			u, err := url.Parse(step.Target)
			if err != nil {
				log.Printf("Error constructing escalation plan for %s in step %d: %s", name, n, err)
				continue STEPS
			}
			s := &notificationStep{
				target: u,
			}
			es.steps = []NotificationSteper{s}
		// NotifyEmail will send an email to the given address.
		case NotifyEmail:
			u, err := url.Parse("mailto://" + step.Target)
			if err != nil {
				log.Printf("Error constructing escalation plan for %s in step %d: %s", name, n, err)
				continue STEPS
			}
			s := &notificationStep{
				target: u,
			}
			es.steps = []NotificationSteper{s}
		default:
			log.Println("Unknown escalation method:", step.Method)
			continue STEPS
		}
		escSteps = append(escSteps, es)
	}
	return escSteps, nil
}

// GetNotificationForTeam will creates a new notificationJob given a team, subject and message.
// A Notification for a team will contain an escalation plan containing multiple escalation steps
// acording to this teams escalation plan. Each step may constist of one or more notification steps,
// depending on the notification method.
func (m *Model) GetNotificationForTeam(name, subject, message string) (notificationJob, error) {
	escSteps, err := m.getEscalationSteps(name)
	if err != nil {
		return notificationJob{}, err
	}
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
