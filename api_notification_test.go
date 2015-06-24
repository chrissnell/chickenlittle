package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/notification"
)

const testCreateNotificationJson = `
{
  "content": "Hello World",
}
`

func TestNotification(t *testing.T) {
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

	// The notification engine must be running, or we'll run into an deadlock
	c.Notify = notification.New(c.Config)

	// prepare the API router
	router := apiRouter()

	// Create a person to test with
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

	// Create a notification plan
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationPlanJson)
	r, err = http.NewRequest("POST", "http://localhost/plan/lancelot", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("CreateNotificaitonPlan request failed")
	}

	// Test NotifyPerson: POST /people/lancelot/notify
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationJson)
	r, err = http.NewRequest("POST", "http://localhost/people/lancelot/notify", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("NotifyPerson request failed")
	}

	// decode the response to extract the UUID for the following API calls
	resp := &NotifyPersonResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}
	uuid := resp.UUID
	if resp.Error != "" {
		t.Fatalf("Notification Error: %s", resp.Error)
	}
	t.Logf("UUID: %s", uuid)
	// we need to give the NotificationEndine some time to pick up the notification job
	time.Sleep(time.Millisecond)

	// Test StopNotification: /notificaionts/{{uuid}}
	w = httptest.NewRecorder()
	r, err = http.NewRequest("DELETE", "http://localhost/notifications/"+uuid, nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("StopNotification request failed: %d", w.Code)
	}

	// Test NotifyPerson: POST /people/lancelot/notify
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationJson)
	r, err = http.NewRequest("POST", "http://localhost/people/lancelot/notify", p)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("NotifyPerson request failed")
	}

	// decode the response to extract the UUID for the following API calls
	resp = &NotifyPersonResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}
	uuid = resp.UUID
	if resp.Error != "" {
		t.Fatalf("Notification Error: %s", resp.Error)
	}
	t.Logf("UUID: %s", uuid)
	// we need to give the NotificationEndine some time to pick up the notification job
	time.Sleep(time.Millisecond)

	// Test StopNotificationClick: ???
	// args: uuid
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://localhost/"+uuid+"/stop", nil)
	if err != nil {
		t.Fatalf("Failed to create new HTTP Request: %s", err)
	}
	clickRouter().ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("StopNotificationClick request failed: %d - %s", w.Code, w.Body)
	}

}
