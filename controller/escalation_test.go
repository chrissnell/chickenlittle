package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreateEscalationPlanJSON = `
{
	"name": "kings",
	"description": "some description",
	"escalation_steps": [
		{
			"timebefore": 1,
			"method": 0,
			"target": ""
		},
		{
			"timebefore": 2,
			"method": 1,
			"target": ""
		},
		{
			"timebefore": 3,
			"method": 2,
			"target": "lancelot"
		},
		{
			"timebefore": 4,
			"method": 3,
			"target": "http://localhost:59876/"
		},
		{
			"timebefore": 5,
			"method": 4,
			"target": "john.doe@localhost"
		}
	]
}
`
const testUpdateEscalationPlanJSON = `
{
	"name": "foobar",
	"description": "foobar",
	"escalation_steps": [
		{
			"timebefore": 900,
			"method": 0
		}
	]
}
`

func TestEscalationPlans(t *testing.T) {
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

	// Test CreateEscalationPlan
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateEscalationPlanJSON)
	r, err = http.NewRequest("POST", "http://localhost/escalation", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreateEscalationPlan request failed: %d", w.Code)
	}

	// Test ListescalationPolicies: GET /escalation
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/escalation", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code == 200 {
		t.Errorf("There should be no escalation policy listing")
	}

	// Test ShowEscalationPlan: GET /escalation/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/escalation/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowEscalationPlan request failed: %d", w.Code)
	}

	// Test UpdateEscalationPlan: PUT /escalation/kings
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testUpdateEscalationPlanJSON)
	r, err = http.NewRequest("PUT", "http://localhost/escalation/kings", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("UpdateEscalationPolicy request failed: %d", w.Code)
	}

	// Test DeleteEscalationPlan: DELETE /escalation/kings
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/escalation/kings", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("DeleteEscalationPlan request failed: %d", w.Code)
	}
}
