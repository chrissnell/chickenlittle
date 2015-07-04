package model

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/db"
)

func TestRotation(t *testing.T) {
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

	// create a model
	m := New(db)

	id := "testpolicy"

	// create a new rotation policy
	rp := RotationPolicy{
		Name:              id,
		Description:       "Example rotation policy",
		RotationFrequency: time.Second,
		RotateTime:        time.Now(),
	}

	// store the RP
	err := m.StoreRotationPolicy(&rp)
	if err != nil {
		t.Errorf("Failed to store the RotationPolicy: %s", err)
	}

	// retrieve the RP
	rp2, err := m.GetRotationPolicy(id)
	if err != nil {
		t.Fatalf("Failed to retrieve the RotationPolicy: %s", err)
	}

	// compare the retrieved RP to the original one
	if rp != *rp2 {
		t.Errorf("The retrieved RotationPolicy should match the original one")
	}

	// retrive using the list method
	rps, err := m.GetAllRotationPolicies()
	if err != nil {
		t.Fatalf("Failed to retrieve the RotationPolicies: %s", err)
	}

	if *rps[0] != rp {
		t.Errorf("The retrieved RotationsPolicies should match the stored ones")
	}

	// delete the RP
	err = m.DeleteRotationPolicy(id)
	if err != nil {
		t.Errorf("Failed to delete the RotationPolicy: %s", err)
	}
}
