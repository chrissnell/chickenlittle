package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreateTeamJSON = `
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

const testUpdateTeamJSON = `
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

// TestTeam will create basic CRUD functionality of a team
func TestTeam(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var p *bytes.Buffer
	var err error

	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the API router
	router := cl.A.APIRouter()

	// Test CreateTeam: POST /teams
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateTeamJSON)
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
	p = bytes.NewBufferString(testUpdateTeamJSON)
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
