package config

import (
	"io/ioutil"
	"os"
	"testing"
)

const testConfig = `
---
service:
  api_listen_address: '127.0.0.2:1234'
  click_listen_address: '127.0.0.3:1234'
  click_url_base: 'http://127.0.0.3:1234/'
  callback_listen_address: '127.0.0.4:1234'
  callback_url_base: 'http://127.0.0.4:1234/'
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
    api_base_url: 'http://localhost:1234'
  mailgun:
    enabled: true
    api_key: 987654
    hostname: localhost
  smtp:
    hostname: localhost
    port: 125
    login: user
    password: pass
    sender: 'mailer@example.com'
`

// TestConfig ensures that we can unmarshall the example config.
func TestConfig(t *testing.T) {
	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	defer func() {
		// remove tempdir
		_ = os.RemoveAll(tempdir)
	}()
	cfgfile := tempdir + "/config.yaml"
	_ = ioutil.WriteFile(cfgfile, []byte(testConfig), 0644)

	c, err := New(cfgfile)
	if err != nil {
		t.Fatalf("Failed to parse config: %s", err)
	}
	if c.Service.APIListenAddr != "127.0.0.2:1234" {
		t.Errorf("Service.APIListenAddr is wrong")
	}
	if c.Service.ClickListenAddr != "127.0.0.3:1234" {
		t.Errorf("Service.ClickListenAddr is wrong")
	}
	if c.Service.ClickURLBase != "http://127.0.0.3:1234/" {
		t.Errorf("Service.ClickURLBase is wrong")
	}
	if c.Service.CallbackListenAddr != "127.0.0.4:1234" {
		t.Errorf("Service.CallbackListenAddr is wrong")
	}
	if c.Service.CallbackURLBase != "http://127.0.0.4:1234/" {
		t.Errorf("Service.CallbackURLBase is wrong")
	}
	if c.Service.DBFile != "/tmp/db" {
		t.Errorf("Service.DBFile is wrong")
	}
	if c.Integrations.HipChat.HipChatAuthToken != "abc" {
		t.Errorf("HipChatAuthToken is wrong")
	}
	if c.Integrations.HipChat.HipChatAnnounceRoom != "def" {
		t.Errorf("HipChatAnnounceRoom is wrong")
	}
	if c.Integrations.VictorOps.APIKey != "hjk" {
		t.Errorf("VictorOps.APIKey is wrong")
	}
	if c.Integrations.Twilio.AccountSID != "12345" {
		t.Errorf("Twilio.AccountSID is wrong")
	}
	if c.Integrations.Twilio.AuthToken != "ABC" {
		t.Errorf("Twilio.AuthToken is wrong")
	}
	if c.Integrations.Twilio.CallFromNumber != "5554321" {
		t.Errorf("Twilio.CallFromNumber is wrong")
	}
	if c.Integrations.Twilio.APIBaseURL != "http://localhost:1234" {
		t.Errorf("Twilio.APIBaseURL is wrong")
	}
	if !c.Integrations.Mailgun.Enabled {
		t.Errorf("Mailgun should be enabled")
	}
	if c.Integrations.Mailgun.APIKey != "987654" {
		t.Errorf("Mailgun.APIKey is wrong")
	}
	if c.Integrations.Mailgun.Hostname != "localhost" {
		t.Errorf("Mailgun.Hostname is wrong")
	}
	if c.Integrations.SMTP.Hostname != "localhost" {
		t.Errorf("SMTP.Hostname is wrong")
	}
	if c.Integrations.SMTP.Port != 125 {
		t.Errorf("SMTP.Port is wrong")
	}
	if c.Integrations.SMTP.Login != "user" {
		t.Errorf("SMTP.Login is wrong")
	}
	if c.Integrations.SMTP.Password != "pass" {
		t.Errorf("SMTP.Password is wrong")
	}
	if c.Integrations.SMTP.Sender != "mailer@example.com" {
		t.Errorf("SMTP.Sender is wrong")
	}
}
