package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/db"
)

func makeTestNotificationPlan(numSteps int) NotificationPlan {
	steps := make([]NotificationStep, 0, numSteps)
	for i := 0; i < numSteps; i++ {
		step := NotificationStep{
			Method:            "NotifyOnDuty",
			NotifyEveryPeriod: time.Second,
			NotifyUntilPeriod: time.Second,
		}
		steps = append(steps, step)
	}
	np := NotificationPlan{
		Username: fmt.Sprintf("testplan%d", numSteps),
		Steps:    steps,
	}
	return np
}

func TestMarshalNotificationPlan(t *testing.T) {
	np := makeTestNotificationPlan(1)
	jb, err := np.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal NotificationPlan: %s", err)
	}
	np2 := NotificationPlan{}
	err = np2.Unmarshal(string(jb))
	if err != nil {
		t.Fatalf("Failed to unmarshal NotificationPlan: %s", err)
	}
	if np.Username != np2.Username || len(np.Steps) != len(np2.Steps) {
		t.Errorf("Unmarshaled Plan should match the original struct")
	}
}

func TestNotificationPlan(t *testing.T) {
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
		np := makeTestNotificationPlan(i)

		err := m.StoreNotificationPlan(&np)
		if err != nil {
			t.Fatalf("Failed to store the escalation plan")
		}
	}

	ep, err := m.GetNotificationPlan("testplan5")
	if err != nil {
		t.Fatalf("Failed to retreive escalation plan: %s", err)
	}
	if len(ep.Steps) != 5 {
		t.Errorf("NotificationPlan 5 should have 5 steps")
	}

	err = m.DeleteNotificationPlan("testplan5")
	if err != nil {
		t.Fatalf("Failed to delete escalation plan: %s", err)
	}

	ep, err = m.GetNotificationPlan("testplan5")
	if err == nil {
		t.Errorf("Should not be able to retrieve deleted plan")
	}
}
