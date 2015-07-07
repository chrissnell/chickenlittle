package model

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrissnell/chickenlittle/db"
)

func TestPeople(t *testing.T) {
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

	// create the model
	m := New(db)

	// create a person to test with
	p := Person{
		Username:            "tester",
		FullName:            "John Doe",
		VictorOpsRoutingKey: "abc",
	}
	// store the person
	err := m.StorePerson(&p)
	if err != nil {
		t.Errorf("Failed to store person: %s", err)
	}
	// retrieve the person from the DB
	p2, err := m.GetPerson("tester")
	if err != nil {
		t.Fatalf("Failed to retrieve person: %s", err)
	}
	// is it the same struct as we put into the DB?
	if p != *p2 {
		t.Errorf("Restored person does not match stored one")
	}
	// get a list of all people to see if the stored one is among them
	peeps, err := m.GetAllPeople()
	if err != nil {
		t.Fatalf("Failed to retrieve all people: %s", err)
	}
	// there should be exaclty one match ...
	if *peeps[0] != p {
		t.Errorf("Restored persons do not match the stored ones")
	}
	// try to delete an non-existing person. It should return an error and not delete anything.
	err = m.DeletePerson("tester-foo")
	if err != nil {
		t.Errorf("Should not be able to delete a non-existing person")
	}
	// fetch all people to see if the numbers are right
	peeps, err = m.GetAllPeople()
	if err != nil {
		t.Fatalf("Failed to retrieve all people: %s", err)
	}
	if len(peeps) != 1 {
		t.Errorf("Should have exactly one person stored")
	}
	// try to delete an existing person
	err = m.DeletePerson("tester")
	if err != nil {
		t.Errorf("Should be able to delete an existing person")
	}
	// the list of people should now be empty
	peeps, err = m.GetAllPeople()
	if err == nil {
		t.Errorf("Should not be able to retrieve an empty list from the DB")
	}

}
