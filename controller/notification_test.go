package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testCreateNotificationJSON = `
{
  "content": "Hello World"
}
`

// TestNotification will test notification of a single person
func TestPersonNotification(t *testing.T) {
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

	// Create a person to test with
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

	// Create a notification plan
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationPlanJSON)
	r, err = http.NewRequest("POST", "http://localhost/plan", p)
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
	p = bytes.NewBufferString(testCreateNotificationJSON)
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
		t.Errorf("Notification Error: %s", resp.Error)
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

	// now enqueue another notification for the same person. this should get another uuid.
	w = httptest.NewRecorder()
	p = bytes.NewBufferString(testCreateNotificationJSON)
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
	cl.A.ClickRouter().ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Fatalf("StopNotificationClick request failed: %d - %s", w.Code, w.Body)
	}
}
