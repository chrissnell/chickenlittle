package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testCreateNotificationPlanJson = `
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

const testUpdateNotificationPlanJson = `
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

	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		// remove tempdir on exit
		_ = os.RemoveAll(tempdir)
	}()
	dbfile := tempdir + "/db"

	// open BoldDB handle
	c.DB.Open(dbfile)
	defer c.DB.Close()

	// prepare the API router
	router := apiRouter()

	// We need a Person to test the notification plans
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreatePersonJson)
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
	p = bytes.NewBufferString(testCreateNotificationPlanJson)
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
	p = bytes.NewBufferString(testUpdateNotificationPlanJson)
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
