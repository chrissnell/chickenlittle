package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

const testConfigYAML = `
---
service:
  api_listen_address: '127.0.0.1:51234'
  click_listen_address: '127.0.0.1:51235'
  click_url_base: 'http://127.0.0.1:51235/'
  callback_listen_address: '127.0.0.1:51236'
  callback_url_base: 'http://127.0.0.1:51236/'
  db_file: '/tmp/db'

integrations:
  hipchat:
    hipchat_auth_token: abc
    hipchat_announce_room: def
  victorops:
    api_key: hjk
  twilio:
    account_sid: 12345
    auth_token: ABC
    call_from_number: 5554321
    api_base_url: 'http://127.0.0.1:51237'
  mailgun:
    enabled: true
    api_key: 987654
    hostname: localhost
  smtp:
    hostname: localhost
    port: 51238
    login: user
    password: pass
    sender: 'mailer@example.com'
`

func TestChickenLittle(t *testing.T) {
	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()
	// create config file
	cfgfile := tempdir + "/config"
	_ = ioutil.WriteFile(cfgfile, []byte(testConfigYAML), 0644)

	cl := New(cfgfile)
	defer cl.Close()
	go cl.Listen()

	// give the http handlers some time to start
	time.Sleep(time.Millisecond)

	// Test REST API Endpoint availability
	r, err := http.Get("http://127.0.0.1:51234/")
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Body.Close()
	if r.StatusCode != 404 {
		t.Errorf("Server should serve no dirindex (yet)")
	}

	// Test Click Endpoint availability
	r, err = http.Get("http://127.0.0.1:51235/")
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Body.Close()
	if r.StatusCode != 404 {
		t.Errorf("Server should serve no dirindex (yet)")
	}

	// Test Callback Endpoint availability
	r, err = http.Get("http://127.0.0.1:51236/")
	if err != nil {
		t.Fatalf("Failed to create HTTP Request: %s", err)
	}
	r.Body.Close()
	if r.StatusCode != 404 {
		t.Errorf("Server should serve no dirindex (yet)")
	}
}
