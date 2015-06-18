package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreatePersonJSON = `
{
  "username": "lancelot",
  "fullname": "Sir Lancelot"
}
`

const testUpdatePersonJSON = `
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

	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the API router
	router := cl.A.APIRouter()

	// Test CreatePerson: POST /people
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreatePersonJSON)
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
	p = bytes.NewBufferString(testUpdatePersonJSON)
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
