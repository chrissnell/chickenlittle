package model

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrissnell/chickenlittle/db"
)

func TestModel(t *testing.T) {
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

	// test a person struct
	p := Person{
		Username: "tester",
		FullName: "John Doe",
	}
	err := m.StorePerson(&p)
	if err != nil {
		t.Errorf("Failed to store person: %s", err)
	}
	p2, err := m.GetPerson("tester")
	if err != nil {
		t.Fatalf("Failed to retrieve person %s", "tester")
	}
	if p != *p2 {
		t.Errorf("Retrieved person should match stored person")
	}
}
