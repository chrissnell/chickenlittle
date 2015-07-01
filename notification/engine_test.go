package notification

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/chrissnell/chickenlittle/config"
	"github.com/chrissnell/chickenlittle/model"
)

type note struct {
	uuid string
	subj string
	mesg string
	meth string
	freq time.Duration
	untl time.Duration
	stps int
	stop chan struct{}
}
type nstep string
type estep string

func (s nstep) Frequency() time.Duration {
	return time.Second
}

func (s nstep) Until() time.Duration {
	return time.Second
}

func (s nstep) NotifyMethod() string {
	return string(s)
}

func (s nstep) Target() *url.URL {
	u, _ := url.Parse(string(s))
	return u
}

func (e estep) ID() string {
	return string(e)
}

func (e estep) Steps() []model.NotificationSteper {
	ns := nstep(e)
	return []model.NotificationSteper{ns, ns}
}

func (n note) Steps() []model.EscalationSteper {
	escSteps := make([]model.EscalationSteper, 0, n.stps)
	for i := 0; i < n.stps; i++ {
		escSteps = append(escSteps, estep(n.meth))
	}
	return escSteps
}
func (n note) ID() string {
	return string(n.uuid)
}
func (n note) Message() string {
	return string(n.mesg)
}
func (n note) Subject() string {
	return string(n.subj)
}
func (n note) Stopper() chan struct{} {
	return n.stop
}

func TestEngine(t *testing.T) {
	notifies := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notifies++
	}))
	defer ts.Close()

	c := config.Config{}
	ng := New(c)
	stopChan := make(chan struct{})
	escSteps := 5
	notifySteps := escSteps * 2
	nr := &note{
		uuid: "1234",
		subj: "Subject",
		mesg: "Message",
		meth: ts.URL,
		freq: 100 * time.Millisecond,
		untl: 500 * time.Millisecond,
		stps: escSteps,
		stop: stopChan,
	}
	ng.EnqueueNotification(nr)
	time.Sleep(time.Second)
	ng.CancelNotification(nr.ID())
	time.Sleep(time.Second)

	// TODO this needs some more finetuning, notifies should equal notifySteps
	if notifies < 1 {
		t.Errorf("Expected %d request, but got %d", notifySteps, notifies)
	}
}
