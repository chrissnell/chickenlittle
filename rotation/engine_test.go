package rotation

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/db"
	"github.com/chrissnell/chickenlittle/model"
)

type testCL struct {
	tempdir string
	c       config.Config
	d       *db.DB
	m       *model.Model
	e       *Engine
}

func newTestCL() *testCL {
	c := &testCL{}
	// create tempdir for fs based tests
	c.tempdir, _ = ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")

	// open config
	c.c = config.NewDefault()

	// open DB
	c.c.Service.DBFile = c.tempdir + "/db"
	c.d = db.New(c.c.Service.DBFile)

	// create model instance
	c.m = model.New(c.d)

	// create a new rotation engine
	c.e = New(c.m)

	return c
}

func (c *testCL) Populate(name string) error {
	// create a team
	initialMembers := []string{"one", "two", "three", "four"}
	team := model.Team{
		Name:           name,
		Members:        initialMembers,
		RotationPolicy: name,
		EscalationPlan: "none",
	}
	err := c.m.StoreTeam(&team)
	if err != nil {
		return err
	}
	return nil
}

func (c *testCL) Close() {
	c.d.Close()
	_ = os.RemoveAll(c.tempdir)
}

func TestPostiveRotationFrequency(t *testing.T) {
	cl := newTestCL()
	defer cl.Close()

	teamName := "tester"
	if err := cl.Populate(teamName); err != nil {
		t.Fatalf("failed to create team %s: %s", teamName, err)
	}
	// create a rotation policy
	rp := model.RotationPolicy{
		Name:              teamName,
		RotationFrequency: time.Second,
		RotateTime:        time.Now(),
	}
	err := cl.m.StoreRotationPolicy(&rp)
	if err != nil {
		t.Fatalf("failed to store rotation policy: %s", err)
	}
	cl.e.UpdatePolicy(teamName)

	// avoid hitting edge cases
	time.Sleep(500 * time.Millisecond)

	// wait for two rotations
	time.Sleep(2 * time.Second)
	// verify the rotated team members
	t2, err := cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "four" {
		t.Errorf("Slot 0 should be four not %s", t2.Members[0])
	}
	if t2.Members[1] != "one" {
		t.Errorf("Slot 1 should be one not %s", t2.Members[1])
	}
	if t2.Members[2] != "two" {
		t.Errorf("Slot 2 should be two not %s", t2.Members[2])
	}
	if t2.Members[3] != "three" {
		t.Errorf("Slot 3 should be three not %s", t2.Members[3])
	}
}

func TestZeroRotationFrequency(t *testing.T) {
	cl := newTestCL()
	defer cl.Close()

	teamName := "tester"
	if err := cl.Populate(teamName); err != nil {
		t.Fatalf("failed to create team %s: %s", teamName, err)
	}
	// create a rotation policy
	rp := model.RotationPolicy{
		Name:              teamName,
		RotationFrequency: 0,
		RotateTime:        time.Now(),
	}
	err := cl.m.StoreRotationPolicy(&rp)
	if err != nil {
		t.Fatalf("failed to store rotation policy: %s", err)
	}
	cl.e.UpdatePolicy("tester")

	// avoid hitting edge cases
	time.Sleep(500 * time.Millisecond)

	// wait some time
	time.Sleep(2 * time.Second)

	// verify the (not!) rotated team members
	t2, err := cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "one" {
		t.Errorf("Slot 0 should be one not %s", t2.Members[0])
	}
	if t2.Members[1] != "two" {
		t.Errorf("Slot 1 should be two not %s", t2.Members[1])
	}
	if t2.Members[2] != "three" {
		t.Errorf("Slot 2 should be three not %s", t2.Members[2])
	}
	if t2.Members[3] != "four" {
		t.Errorf("Slot 3 should be four not %s", t2.Members[3])
	}
}

func TestRotateTimeFuture(t *testing.T) {
	cl := newTestCL()
	defer cl.Close()

	teamName := "tester"
	if err := cl.Populate(teamName); err != nil {
		t.Fatalf("failed to create team %s: %s", teamName, err)
	}
	// create a rotation policy
	rp := model.RotationPolicy{
		Name:              teamName,
		RotationFrequency: time.Second,
		RotateTime:        time.Now().Add(2 * time.Second),
	}
	err := cl.m.StoreRotationPolicy(&rp)
	if err != nil {
		t.Fatalf("failed to store rotation policy: %s", err)
	}
	cl.e.UpdatePolicy(teamName)

	// avoid hitting edge cases
	time.Sleep(500 * time.Millisecond)

	time.Sleep(time.Second)

	// verify the team order before rotation
	t2, err := cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "one" {
		t.Errorf("Slot 0 should be one not %s before the rotate start time", t2.Members[0])
	}

	// wait until the regular rotations should begin
	time.Sleep(time.Second)
	// verify the team order after the first rotation
	t2, err = cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "two" {
		t.Errorf("Slot 0 should be two not %s", t2.Members[0])
	}
}

func TestRotateTimePast(t *testing.T) {
	cl := newTestCL()
	defer cl.Close()

	teamName := "tester"
	if err := cl.Populate(teamName); err != nil {
		t.Fatalf("failed to create team %s: %s", teamName, err)
	}
	// create a rotation policy
	rp := model.RotationPolicy{
		Name:              teamName,
		RotationFrequency: 2 * time.Second,
		RotateTime:        time.Now().Add(-3 * time.Second),
	}
	// x - past rotations
	// X - future rotations
	// 0 - Now()
	// -3 -2 -1  0  1  2  3  4
	//  x     x     X     X
	err := cl.m.StoreRotationPolicy(&rp)
	if err != nil {
		t.Fatalf("failed to store rotation policy: %s", err)
	}
	cl.e.UpdatePolicy(teamName)

	// avoid hitting edge cases
	time.Sleep(500 * time.Millisecond)

	// verify the team order before rotation
	t2, err := cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "one" {
		t.Errorf("Slot 0 should be one not %s before the rotate start time", t2.Members[0])
	}

	// wait until after the first regular rotation
	time.Sleep(2 * time.Second)
	// verify the team order before rotation
	t2, err = cl.m.GetTeam(teamName)
	if err != nil {
		t.Fatalf("failed to retrieve team: %s", err)
	}
	if t2.Members[0] != "two" {
		t.Errorf("Slot 0 should be two not %s after the rotate start time", t2.Members[0])
	}
}
