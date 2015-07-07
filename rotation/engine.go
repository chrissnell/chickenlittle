package rotation

import (
	"log"
	"sync"

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
