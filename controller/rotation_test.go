package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreateRotationPolicyJSON = `
{
	"name": "kings",
	"description": "some description",
	"frequency": 0,
	"time": "2015-07-06T00:00:01Z"
}
`
const testUpdateRotationPolicyJSON = `
{
	"name": "foobar",
	"description": "foobar",
	"frequency": 1,
	"time": "2015-07-06T00:00:01Z"
}
`

func TestRotationPolicies(t *testing.T) {
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

	// Test CreateRotationPolicy
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateRotationPolicyJSON)
	r, err = http.NewRequest("POST", "http://localhost/rotation", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreateRotationPolicy request failed: %d", w.Code)
	}

	// Test ListRotationPolicies: GET /rotation
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/rotation", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code == 200 {
		t.Errorf("There should be no rotation policy listing")
	}

	// Test ShowRotationPolicy: GET /rotation/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/rotation/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowRotationPolicy request failed: %d", w.Code)
	}

	// Test UpdateRotationPolicy: PUT /rotation/kings
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testUpdateRotationPolicyJSON)
	r, err = http.NewRequest("PUT", "http://localhost/rotation/kings", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("UpdateEscalationPolicy request failed: %d", w.Code)
	}

	// Test DeleteRotationPolicy: DELETE /rotation/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/rotation/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("DeleteRotationPolicy request failed: %d", w.Code)
	}
}
