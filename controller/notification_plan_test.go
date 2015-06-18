package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCreateNotificationPlanJSON = `
[
  {
    "method": "noop://2108675309",
    "notify_every_period": 0,
    "notify_until_period": 300000000000
  },
  {
    "method": "noop://2105551212",
    "notify_every_period": 900000000000,
    "notify_until_period": 0
  }
]
`

const testUpdateNotificationPlanJSON = `
[
  {
    "method": "noop://2108675309",
    "notify_every_period": 0,
    "notify_until_period": 300000000000
  }
]
`

func TestNotificationPlan(t *testing.T) {
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
		t.Fatalf("CreatePerson request failed")
	}

	// Test CreateNotificationPlan: POST /plan/{{username}}
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationPlanJSON)
	r, err = http.NewRequest("POST", "http://localhost/plan/lancelot", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("CreateNotificationPlan request failed")
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
		t.Errorf("ShowNotificationPlan request failed")
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
		t.Errorf("UpdateNotificationPlan request failed")
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
		t.Errorf("DeleteNotificationPlan request failed")
	}

}
