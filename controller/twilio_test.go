package controller

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/twinj/uuid"
)

func TestReceiveSMSReply(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var err error

	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the callback router
	router := cl.api.CallbackRouter()
	// POST /sms
	w = httptest.NewRecorder()
	data := url.Values{
		"From": {"10015554321"},
		"Body": {"123"},
	}
	r, err = http.NewRequest("POST", "http://localhost/sms", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 404 {
		t.Errorf("Expected code %d but got %d", 404, w.Code)
	}

	// TODO now we've confirmed that unknown replies are ignored test that we recognize known ones

	// noop
	t.Logf("Not yet (fully) implemented")
}

func TestReceiveCallback(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var err error

	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the callback router
	router := cl.api.CallbackRouter()

	// TODO we should create a notification first to have an valid uuid
	uuid := uuid.NewV4().String()

	// POST /{{uuid}}/callback
	w = httptest.NewRecorder()
	data := url.Values{}
	r, err = http.NewRequest("POST", "http://localhost/"+uuid+"/callback", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 200 {
		t.Errorf("Expected code %d but got %d", 404, w.Code)
	}
	// TODO should decode the CallbackResponse and compare Message and UUID
}

func TestReceiveDigits(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	//var p *bytes.Buffer
	var err error

	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the callback router
	router := cl.api.CallbackRouter()

	// TODO we should create a notification first to have an valid uuid
	uuid := uuid.NewV4().String()

	// TODO we should provide an httptest Server to receive and check the Twilio callback

	// POST /{{uuid}}/digits
	w = httptest.NewRecorder()
	data := url.Values{
		"Digits":  {"123"},
		"CallSid": {"456"},
	}
	r, err = http.NewRequest("POST", "http://localhost/"+uuid+"/digits", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 404 {
		t.Errorf("Expected code %d but got %d", 404, w.Code)
	}

	// noop
	t.Logf("Not yet (fully) implemented")
}

func TestGenerateTwiML(t *testing.T) {
	var w *httptest.ResponseRecorder
	var r *http.Request
	var err error

	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// prepare the callback router
	router := cl.api.CallbackRouter()

	// TODO we should create a notification first to have an valid uuid
	uuid := uuid.NewV4().String()

	// TODO we should provide an httptest Server to receive and check the Twilio callback

	// POST /{{uuid}}/twiml/{action}
	w = httptest.NewRecorder()
	r, err = http.NewRequest("POST", "http://localhost/"+uuid+"/twiml/notify", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	router.ServeHTTP(w, r)
	// verify response
	if w.Code != 404 {
		t.Errorf("Expected code %d but got %d", 404, w.Code)
	}

	// noop
	t.Logf("Not yet (fully) implemented")
}
