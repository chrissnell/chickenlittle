package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreateNotificationPlanJSON = `
{
	"username": "lancelot",
	"steps":
		[
		  {
		    "method": "noop://2108675309",
		    "notify_every_period": 1000000000,
		    "notify_until_period": 3000000000
		  },
		  {
		    "method": "noop://2105551212",
		    "notify_every_period": 1000000000,
		    "notify_until_period": 3000000000
		  }
		]
}
`

const testUpdateNotificationPlanJSON = `
{
	"username": "foobar",
	"steps":
		[
		  {
		    "method": "noop://2108675309",
		    "notify_every_period": 1000000000,
		    "notify_until_period": 3000000000
		  }
		]
}
`

// TestNotificationPlan will test the basic CRUD functionality for an notification plan
func TestNotificationPlan(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var p *bytes.Buffer
	var err error

	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create new test client: %s", err)
	}
	defer cl.Close()

	// prepare the API router
	router := cl.api.APIRouter()

	// We need a Person to test the notification plans
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreatePersonJSON)
	r, err = http.NewRequest("POST", "http://localhost/people", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("CreatePerson request failed: %d", w.Code)
	}

	// Test CreateNotificationPlan: POST /plan/{{username}}
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationPlanJSON)
	r, err = http.NewRequest("POST", "http://localhost/plan", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreateNotificationPlan request failed: %d", w.Code)
	}

	// Test ShowNotificationPlan: GET /plan/lancelot
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/plan/lancelot", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("ShowNotificationPlan request failed: %d", w.Code)
	}

	// Test UpdateNotificaitonPlan: PUT /plan/lancelot
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testUpdateNotificationPlanJSON)
	r, err = http.NewRequest("PUT", "http://localhost/plan/lancelot", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("UpdateNotificationPlan request failed: %d", w.Code)
	}

	// Test DeleteNotificationPlan: DELETE /plan/lancelot
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/plan/lancelot", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("DeleteNotificationPlan request failed: %d", w.Code)
	}

}
