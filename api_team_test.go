package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testCreateTeamJson = `
{
  "description": "Kings", 
  "members": [
	  "lancelot",
    "arthur"
  ], 
  "name": "kings", 
  "rotation": {
    "description": "none", 
    "frequency": 0
  }, 
  "steps": [
		{
  	  "method": 0, 
  	  "target": "", 
  	  "timebefore": 3600
  	}
	]
}
`

const testUpdateTeamJson = `
{
  "description": "Kings", 
  "members": [
    "arthur"
  ], 
  "name": "kings", 
  "rotation": {
    "description": "none", 
    "frequency": 0
  }, 
  "steps": [
		{
  	  "method": 0, 
  	  "target": "", 
  	  "timebefore": 3600
  	}
	]
}
`

func TestTeam(t *testing.T) {
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

	// Test CreateTeam: POST /teams
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateTeamJson)
	r, err = http.NewRequest("POST", "http://localhost/teams", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreateTeams request failed")
	}

	// Test ListTeams: GET /teams
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/teams", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ListTeams request failed")
	}

	// Test ShowTeam: GET /teams/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/teams/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowTeam request failed")
	}

	// Test UpdateTeam: PUT /teams/kings
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testUpdateTeamJson)
	r, err = http.NewRequest("PUT", "http://localhost/teams/kings", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("UpdateTeam request failed")
	}

	// Test DeleteTeam: DELETE /teams/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/teams/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("DeleteTeam request failed")
	}

}
