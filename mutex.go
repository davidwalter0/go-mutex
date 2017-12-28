// Note that in the following the _scope_ closure for a deferred call
// of a lock is the function exit, therefore scoped lock calls lock at
// the defer statement and close at the end of the function not at the
// block *{}* level

package mutex

import (
	"fmt"
	"sync"
)

// debug var for testing
var testDebug bool

func scopedTrace(text string) func() {
	fmt.Println(text)
	return func() {
		fmt.Println(text)
	}
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", 2*n, " ")
}

// Mutex local synonym for sync.Mutex for receiver methods
type Mutex sync.Mutex

// NewMutex return a mutex pointer
func NewMutex() *Mutex {
	return &Mutex{}
}

// NewMonitor return a deferable preinitialized private mutex closure
// monitor function.
// defer scope: acquire lock on entry, and release on scope closure
// created anonymously in the function closure

func NewMonitor() func(...interface{}) func() {
	var mutex = NewMutex()
	return func(args ...interface{}) func() {
		var i int
		if len(args) > 0 {
			switch args[0].(type) {
			case int:
				i = args[0].(int)
			}
		}
		if testDebug {
			defer scopedTrace(spaces(i) + fmt.Sprintf("monitor lock %d", i))()
		}
		mutex.Lock()
		return func() {
			if testDebug {
				defer scopedTrace(spaces(i) + fmt.Sprintf("monitor unlock %d", i))()
			}
			mutex.Unlock()
		}
	}
}

// Lock the mutex
func (mutex *Mutex) Lock() {
	(*sync.Mutex)(mutex).Lock()
}

// Unlock the mutex
func (mutex *Mutex) Unlock() {
	(*sync.Mutex)(mutex).Unlock()
}

// MonitorTrace block scoped mutex with depth print
// defer mutex.MonitorTrace()()
// prefer to use example from tests with defer GuardedTrace()()
func (mutex *Mutex) MonitorTrace(args ...interface{}) func() {
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

// Monitor block scoped mutex block scoped mutex returns function for
// defer call. Ex:
// defer mutex.Monitor()()
func (mutex *Mutex) Monitor() func() {
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

// Guard alias of Monitor() block scoped mutex returns function for
// defer call. Ex:
// defer mutex.Guard()()
func (mutex *Mutex) Guard() func() {
	return mutex.Monitor()
}
