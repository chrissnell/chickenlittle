package db

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		// remove tempdir
		_ = os.RemoveAll(tempdir)
	}()
	dbfile := tempdir + "/db"

	// open BoldDB handle
	db := New(dbfile)
	defer db.Close()

	refBucket := "testbucket"
	refKey := "testkey"
	refVal := "testvalue"

	err := db.Store(refBucket, refKey, refVal)
	if err != nil {
		t.Errorf("Failed to store key: %s", err)
	}

	val, err := db.Fetch(refBucket, refKey)
	if err != nil {
		t.Errorf("Failed to retrieve stored key: %s", err)
	}
	if val != refVal {
		t.Errorf("Value should be %s not %s", refVal, val)
	}
	vals, err := db.FetchAll(refBucket)
	if err != nil {
		t.Errorf("Failed to retrieve stored keys: %s", err)
	}
	if vals[0] != refVal {
		t.Errorf("Value should be %s not %s", refVal, vals[0])
	}
	err = db.Delete(refBucket, refKey)
	if err != nil {
		t.Errorf("Failed to delete a stored key: %s", err)
	}
	val, err = db.Fetch(refBucket, refKey)
	if err == nil {
		t.Errorf("Should not be able to retrieve the deleted key")
	}
}
