// datastruct.go
package main

import (
	"container/list"
)

/*****************************************************************************/

type LockArea struct {
	locks   map[string]*list.List
	clients map[*Client]map[string]bool
}

/*****************************************************************************/

func NewLockArea() *LockArea {
	return &LockArea{
		locks:   make(map[string]*list.List),
		clients: make(map[*Client]map[string]bool),
	}
}

/*****************************************************************************/

func (lo *LockArea) Add(clt *Client, name string) bool {

	if clist, ok := lo.locks[name]; ok {
		if clist.Front().Value.(*Client) != clt {
			clist.PushBack(clt)
			lo.clients[clt][name] = false
			return false
		} else {
			return true
		}
	} else {
		clist = list.New()
		clist.PushBack(clt)
		lo.locks[name] = clist
		lo.clients[clt][name] = true
		return true
	}
}

/*****************************************************************************/

func (lo *LockArea) Remove(clt *Client, name string) (*Client, bool) {

	clist, lok := lo.locks[name]
	if !lok {
		return nil, false
	}
	if present := lo.clients[clt][name]; !present {
		return nil, false
	}
	e := clist.Front()
	c := e.Value.(*Client)
	if c != clt {
		panic("LockArea data structure is corrupted")
	}
	clist.Remove(e)
	delete(lo.clients[clt], name)

	e = clist.Front()
	if e == nil {
		delete(lo.locks, name)
		return nil, true
	} else {
		c = e.Value.(*Client)
		lo.clients[c][name] = true
		return c, true
	}
}

/*****************************************************************************/

func (lo *LockArea) AddClient(clt *Client) {

	lo.clients[clt] = make(map[string]bool)
}

/*****************************************************************************/

func (lo *LockArea) RemoveClient(clt *Client) []*Client {

	res := []*Client{}

	for name, locked := range lo.clients[clt] {
		if locked {
			if next, _ := lo.Remove(clt, name); next != nil {
				res = append(res, next)
			}
		} else {
			l := lo.locks[name]
			for e := l.Front(); e != nil; {
				next := e.Next()
				c := e.Value.(*Client)
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
