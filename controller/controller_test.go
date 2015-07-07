package controller

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/db"
	"github.com/chrissnell/chickenlittle/model"
	"github.com/chrissnell/chickenlittle/notification"
	"github.com/chrissnell/chickenlittle/rotation"
)

type clTemp struct {
	Config  config.Config
	db      *db.DB
	notify  *notification.Engine
	api     *Controller
	model   *model.Model
	tempdir string
}

func newTestClient() (*clTemp, error) {
	var err error
	c := &clTemp{}

	// create tempdir for fs based tests
	c.tempdir, err = ioutil.TempDir(os.TempDir(), "chickenlittle-tests-")
	if err != nil {
		return c, err
	}
	c.Config = config.NewDefault()
	c.Config.Service.DBFile = c.tempdir + "/db"
	c.db = db.New(c.Config.Service.DBFile)
	c.model = model.New(c.db)
	c.notify = notification.New(c.Config)
	r := rotation.New(c.model)
	c.api = New(c.Config, c.model, c.notify, r)

	return c, nil
}

func (c *clTemp) Close() {
	c.db.Close()
	_ = os.RemoveAll(c.tempdir)
}

// TestAPI will only ensure we can create an test client with API, model, db, et. al.
func TestAPI(t *testing.T) {
	cl, err := newTestClient()
	if err != nil {
		t.Fatalf("Failed to create new test client: %s", err)
	}
	defer cl.Close()
}
