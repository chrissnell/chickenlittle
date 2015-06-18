package ne

import (
	"testing"
	"time"
)

type note string
type step string

func (s step) Frequency() time.Duration {
	return time.Second
}

func (s step) Until() time.Duration {
	return time.Second
}

func (s step) NotifyMethod() string {
	return string(s)
}

func (n note) NextStep() (NotificationStep, bool) {
	return step("noop://123456"), false
}
func (n note) ID() string {
	return string(n)
}
func (n note) Message() string {
	return string(n)
}
func (n note) Subject() string {
	return string(n)
}

func TestEngine(t *testing.T) {
	neCfg := Config{}
	ng := New(neCfg)
	nr := note("foo")
	ng.EnqueueNotification(nr)
	time.Sleep(time.Second)
	ng.CancelNotification(nr.ID())
}
