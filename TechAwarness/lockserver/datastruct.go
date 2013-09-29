// datastruct.go
package main

import (
	"container/ring"
)

type LockArea struct {
	locks map[string]Lock
	clients map[*Client]ring.Ring
}

type Lock struct {
	clients *ring.Ring
}
