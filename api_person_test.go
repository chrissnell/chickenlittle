package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testCreatePersonJson = `
{
  "username": "lancelot",
  "fullname": "Sir Lancelot"
}
`

const testUpdatePersonJson = `
{
	"username": "lancelot",
	"fullname": "Sir"
}
`

func TestPerson(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var p *bytes.Buffer
	var err error

	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		// remove tempdir
		_ = os.RemoveAll(tempdir)
	}()
	dbfile := tempdir + "/db"

	// open BoldDB handle
	c.DB.Open(dbfile)
	defer c.DB.Close()

	// prepare the API router
	router := apiRouter()

	// Test CreatePerson: POST /people
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreatePersonJson)
	r, err = http.NewRequest("POST", "http://localhost/people", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreatePerson request failed")
	}

	// Test ListPeople: GET /people
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/people", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ListPeople request failed")
	}

	// Test ShowPerson: GET /people/lancelot
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/people/lancelot", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowPerson request failed")
	}

	// Test UpdatePerson: PUT /people/lancelog
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testUpdatePersonJson)
	r, err = http.NewRequest("PUT", "http://localhost/people/lancelot", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("UpdatePerson request failed")
	}

	// Test DeletePerson: DELETE /people/lancelot
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/people/lancelot", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowPerson request failed")
	}

}
