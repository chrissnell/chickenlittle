package rotation

import (
	"log"
	"sync"
	"time"

	"github.com/chrissnell/chickenlittle/model"
)

// Engine is the rotation engine managing all rotation policy watcher
type Engine struct {
	model   *model.Model
	mutex   *sync.Mutex        // protects all below
	watcher map[string]Watcher // one watcher per rotation policy
}

// New creates a new rotation engine
func New(m *model.Model) *Engine {
	e := &Engine{
		model:   m,
		mutex:   &sync.Mutex{},
		watcher: make(map[string]Watcher),
	}
	return e
}

// UpdatePolicy tells the rotation engine to update or create a rotation watcher for the given policy
func (e *Engine) UpdatePolicy(name string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	log.Println("Updating Policy for", name)
	// if a watcher exists tell it to re-examine his config
	if w, found := e.watcher[name]; found {
		log.Println("Watcher found. Notifying ...")
		w.Notify()
		return
	}

	// if no watcher exists yet, create one
	e.watcher[name] = NewWatcher(e, name)
	log.Println("Watcher not found. Created new one.")
}

/*
Control Flow:

create policy p -> e.UpdatePolicy(p) -> NewWatcher(p) -> run
update policy p	-> e.UpdatePolicy(p) -> w.Notify -> updateChan -> run -> m.Get -> restart timer
delete policy p -> e.UpdatePolicy(p) -> w.Notify -> updateChan -> run -> break -> e.unregister(w) -> return

*/

// unregister will remove one Watcher from the engine
func (e *Engine) unregister(w *Watcher) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.watcher, w.name)
}

// Watcher is a single watcher for a rotation policy
type Watcher struct {
	name       string // name of the watched Team
	e          *Engine
	policy     model.RotationPolicy
	updateChan chan struct{}
}

// NewWatcher will create a new watched tied to the engine
func NewWatcher(e *Engine, name string) Watcher {
	w := Watcher{
		name:       name,
		e:          e,
		updateChan: make(chan struct{}, 10),
	}
	go w.run()
	return w
}

// Notify should be called whenever the underlying primitives (rotation policy, team)
// change.
func (w *Watcher) Notify() {
	w.updateChan <- struct{}{}
}

func (w *Watcher) initTicker() <-chan time.Time {
	// a rotation requency of 0 means: no automatic rotations, so create a dummy
	// channel here that will never receive an tick. In that case we don't even
	// care about an RotateTime
	if w.policy.RotationFrequency < 1 {
		log.Printf("[ %10s ] RotationFrequency is zero. Disabling any rotations.", w.name)
		return make(chan time.Time)
	}

	if w.policy.RotateTime.Before(time.Now()) {
		// not so easy, this policy should have been started in the past, but we have to
		// make sure we get the offsets right so that the next rotation happens at
		// (Now()-RotateTime) % RotationFrequency ... ?!
		sinceFirstRotation := time.Now().Sub(w.policy.RotateTime)
		nextRotation := sinceFirstRotation % w.policy.RotationFrequency
		log.Printf("[ %10s ] RotateTime is in the past. First rotation will be in %fs", w.name, nextRotation.Seconds())
		return time.After(nextRotation)
	}

	// easy, this policy should start at some point in the future
	// so use the ticker to start regular rotations at the point in time
	log.Printf("[ %10s ] RotateTime is in the future. First rotation will be in %fs", w.name, w.policy.RotateTime.Sub(time.Now()).Seconds())
	return time.After(w.policy.RotateTime.Sub(time.Now()))
}

func (w *Watcher) run() {
	var ticker <-chan time.Time
	w.Notify() // initialize ourself inside the regular update loop

	for {
		select {
		case <-w.updateChan:
			log.Printf("[ %10s ] Watcher checking for updates", w.name)
			changed, deleted := w.update()
			if deleted {
				log.Printf("[ %10s ] rotation policy was deleted", w.name)
				break
			}
			if changed {
				log.Printf("[ %10s ] rotation policy was updated", w.name)
				ticker = w.initTicker()
			}
		case <-ticker:
			w.rotate()
			log.Printf("[ %10s ] Watcher rotated team", w.name)
			ticker = time.After(w.policy.RotationFrequency)
		}
	}
	w.e.unregister(w)
	log.Printf("[ %10s ] watcher shuting down", w.name)
}

// update will update our copy of the rotation policy and return if it was changed or deleted
// this method is not thread-safe, as it should only be called from within the (single)
// run loop.
func (w *Watcher) update() (changed bool, deleted bool) {
	newPolicy, err := w.e.model.GetRotationPolicy(w.name)
	if err != nil {
		changed = true
		deleted = true
		return // the policy was deleted
	}
	if *newPolicy != w.policy {
		w.policy = *newPolicy
		changed = true
		deleted = false
		return // the policy was changed
	}
	// no changes to the policy
	return
}

// rotate will do a simple rotation on the team being watched
func (w *Watcher) rotate() {
	team, err := w.e.model.GetTeam(w.name)
	if err != nil {
		log.Printf("failed to look up team %s: %s", w.name, err)
		return
	}

	log.Printf("rotating team: %s", w.name)
	members := make([]string, 0, len(team.Members))
	for i := 1; i < len(team.Members); i++ {
		members = append(members, team.Members[i])
	}
	members = append(members, team.Members[0])

	for n, m := range members {
		log.Printf("%d: %s", n, m)
	}

	team.Members = members
	err = w.e.model.StoreTeam(team)
	if err != nil {
		log.Printf("failed to store the rotated team %s: %s", w.name, err)
	}
	return
}
