package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/db"
)

func makeTestEscalationPlan(numSteps int) EscalationPlan {
	steps := make([]EscalationStep, 0, numSteps)
	for i := 0; i < numSteps; i++ {
		step := EscalationStep{
			TimeBeforeEscalation: time.Second,
			Method:               NotifyOnDuty,
			Target:               fmt.Sprintf("%d", i),
		}
		steps = append(steps, step)
	}
	ep := EscalationPlan{
		Name:        fmt.Sprintf("testplan%d", numSteps),
		Description: fmt.Sprintf("a simple test plan with %d steps", numSteps),
		Steps:       steps,
	}
	return ep
}

func TestMarshalEscalationPlan(t *testing.T) {
	ep := makeTestEscalationPlan(1)
	jb, err := ep.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal EscalationPlan: %s", err)
	}
	ep2 := EscalationPlan{}
	err = ep2.Unmarshal(string(jb))
	if err != nil {
		t.Fatalf("Failed to unmarshal EscalationPlan: %s", err)
	}
	if ep.Name != ep2.Name || ep.Description != ep2.Description || len(ep.Steps) != len(ep2.Steps) {
		t.Errorf("Unmarshaled Plan should match the original struct")
	}
}

func TestEscalationPlan(t *testing.T) {
	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		// remove tempdir
		_ = os.RemoveAll(tempdir)
	}()
	dbfile := tempdir + "/db"

	// open BoldDB handle
	db := db.New(dbfile)
	defer db.Close()

	// create model
	m := New(db)

	// create some plans
	for i := 0; i < 10; i++ {
		ep := makeTestEscalationPlan(i)

		err := m.StoreEscalationPlan(&ep)
		if err != nil {
			t.Fatalf("Failed to store the escalation plan")
		}
	}

	eps, err := m.GetAllEscalationPlans()
	if err != nil {
		t.Fatalf("Failed to retrieve escalation plans: %s", err)
	}
	if len(eps) != 10 {
		t.Errorf("Should have received %d plans not %d", 10, len(eps))
	}

	ep, err := m.GetEscalationPlan("testplan5")
	if err != nil {
		t.Fatalf("Failed to retreive escalation plan: %s", err)
	}
	if len(ep.Steps) != 5 {
		t.Errorf("EscalationPlan 5 should have 5 steps")
	}

	err = m.DeleteEscalationPlan("testplan5")
	if err != nil {
		t.Fatalf("Failed to delete escalation plan: %s", err)
	}

	ep, err = m.GetEscalationPlan("testplan5")
	if err == nil {
		t.Errorf("Should not be able to retrieve deleted plan")
	}
}
