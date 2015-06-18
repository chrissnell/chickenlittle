package model

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrissnell/chickenlittle/db"
)

func TestTeam(t *testing.T) {
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

	// .. and a team to test with
	te := Team{
		Name:        "tester",
		Description: "Just some testers",
		Members: []string{
			"Lancelot",
			"Arthur",
		},
		RotationPolicy: "rp1",
		EscalationPlan: "ep1",
	}

	// store the team
	err := m.StoreTeam(&te)
	if err != nil {
		t.Errorf("Failed to store team: %s", err)
	}

	// retrieve it again
	te2, err := m.GetTeam("tester")
	if err != nil {
		t.Fatalf("Failed to retrieve team: %s", err)
	}
	// TODO compare te and te2, they are not trival comparable
	_ = te2
	t.Logf("Should compare te and te2")

	// get the list of teams
	tes, err := m.GetAllTeams()
	if err != nil {
		t.Fatalf("Failed to retrieve teams: %s", err)
	}
	if len(tes) != 1 {
		t.Errorf("There should be exactly one team")
	}

	// delete the team
	err = m.DeleteTeam("tester")
	if err != nil {
		t.Errorf("Failed to delete team: %s", err)
	}
}
