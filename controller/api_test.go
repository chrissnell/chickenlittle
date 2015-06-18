package controller

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/db"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
)

const testConfigYAML = `
---
service:
  api_listen_address: '127.0.0.1:51234'
  click_listen_address: '127.0.0.1:51234'
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

type TestCL struct {
	Config  config.Config
	D       *db.DB
	N       *notification.Engine
	A       *Controller
	M       *model.Model
	tempdir string
}

// NewTestCL does all the necessary setup for a simple test client
func NewTestCL() (*TestCL, error) {
	// create tempdir for fs based tests
	tempdir, _ := ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	// create config file
	cfgfile := tempdir + "/config"
	_ = ioutil.WriteFile(cfgfile, []byte(testConfigYAML), 0644)

	// open config
	c, err := config.New(cfgfile)
	if err != nil {
		return nil, err
	}

	// open DB
	c.Service.DBFile = tempdir + "/db"
	d := db.New(c.Service.DBFile)

	// create model instance
	m := model.New(d)

	// create notification engine instance
	n := notification.New(c)

	// create API instance
	a := New(c, m, n)

	// create the test client
	cl := &TestCL{
		Config:  c,
		D:       d,
		N:       n,
		A:       a,
		M:       m,
		tempdir: tempdir,
	}
	return cl, nil
}

// Close performs the necessary teardown of the test client.
func (c *TestCL) Close() {
	c.D.Close()
	_ = os.RemoveAll(c.tempdir)
}

func TestAPI(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// this test only ensure we can create an test client with API, model, db, et. al.
}
