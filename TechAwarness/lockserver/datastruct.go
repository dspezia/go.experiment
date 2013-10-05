// This file contains the locking data structure management code.

package lockserver

import "container/list"

/*****************************************************************************/

// Replier represents a client that can process a reply message
type Replier interface {
	Reply(*MessageReply)
}

/*****************************************************************************/

// LockArea is the data structure responsible of tracking who locks what, and
// what is locked by who.
type LockArea struct {
	locks   map[string]*list.List       // Map associating locks to list of clients
	clients map[Replier]map[string]bool // Map associating repliers to map of locks
}

/*****************************************************************************/

// NewLockArea builds a new LockArea object
func NewLockArea() *LockArea {
	return &LockArea{
		locks:   make(map[string]*list.List),
		clients: make(map[Replier]map[string]bool),
	}
}

/*****************************************************************************/

// Add must be called to notify a locking event
func (lo *LockArea) Add(clt Replier, name string) bool {

	// Check if lock already exists
	if clist, ok := lo.locks[name]; ok {
		// Check if the client has already locked the same object
		if clist.Front().Value.(Replier) != clt {
			// No: client is just queued, do no reply
			clist.PushBack(clt)
			lo.clients[clt][name] = false
			return false
		} else {
			// Yes: just ignore and reply
			return true
		}
	} else {
		// Create new lock object and grant it to the client, and reply
		clist = list.New()
		clist.PushBack(clt)
		lo.locks[name] = clist
		lo.clients[clt][name] = true
		return true
	}
}

/*****************************************************************************/

// Remove is called to notify an unlock
func (lo *LockArea) Remove(clt Replier, name string) (Replier, bool) {

	// Sanity check: the lock must exist
	clist, lok := lo.locks[name]
	if !lok {
		return nil, false
	}
	// Sanity check: the client must hold the lock
	if present := lo.clients[clt][name]; !present {
		return nil, false
	}
	e := clist.Front()
	c := e.Value.(Replier)
	if c != clt {
		panic("LockArea data structure is corrupted")
	}

	// Remove clients from lock list, delete lock from client map
	clist.Remove(e)
	delete(lo.clients[clt], name)

	// Check whether the lock can be granted to another client
	e = clist.Front()
	if e == nil {
		// No other lock intent
		delete(lo.locks, name)
		return nil, true
	} else {
		// Found a lock intent: grant the lock
		c = e.Value.(Replier)
		lo.clients[c][name] = true
		return c, true
	}
}

/*****************************************************************************/

// AddClient is called to notify a new client
func (lo *LockArea) AddClient(clt Replier) {

	lo.clients[clt] = make(map[string]bool)
}

/*****************************************************************************/

// Remove client is called to notify a client disconnection.
// It can be due to a normal disconnection, or a crash.
func (lo *LockArea) RemoveClient(clt Replier) []Replier {

	res := []Replier{}

	// Iterate on all the locks related to the client
	for name, locked := range lo.clients[clt] {
		if locked {
			// This lock was granted to the client, it must be released
			if next, _ := lo.Remove(clt, name); next != nil {
				res = append(res, next)
			}
		} else {
			// This was only a lock intent, but it has to be removed
			l := lo.locks[name]
			for e := l.Front(); e != nil; {
				next := e.Next()
				c := e.Value.(Replier)
				if c == clt {
					l.Remove(e)
				}
				e = next
			}
			if l.Front() == nil {
				delete(lo.locks, name)
			}
		}
	}

	delete(lo.clients, clt)
	return res
}

/*****************************************************************************/
