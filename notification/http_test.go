package notification

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chrissnell/chickenlittle/config"
)

func TestHttp(t *testing.T) {
	recvSubject := ""
	recvMessage := ""
	recvUUID := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recvSubject = r.FormValue("subject")
		recvMessage = r.FormValue("message")
		recvUUID = r.FormValue("uuid")
	}))
	defer ts.Close()

	c := config.NewDefault()
	e := New(c)
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("Failed to parse URL: %s", err)
	}
	err = e.CallHTTP(u, "subj", "msg", "uuid")
	if err != nil {
		t.Errorf("Failed to make HTTP request: %s", err)
	}

	if recvSubject != "subj" {
		t.Errorf("Received subject wrong")
	}

	if recvMessage != "msg" {
		t.Errorf("Received message wrong")
	}

	if recvUUID != "uuid" {
		t.Errorf("Received uuid wrong")
	}
}
