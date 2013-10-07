package main

import (
	"fmt"
)

// B1 OMIT
type Rectangle struct{ l, h float64 }

func (self Rectangle) Area() float64 {
	return self.l * self.h
}

func (self *Rectangle) Zero() {
	self.l, self.h = 0, 0
}

// E1 OMIT

// B2 OMIT
type Length float64

func (self Length) Square() float64 {
	return self * self
}

// E2 OMIT

// B3 OMIT
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Empty interface{}

// E3 OMIT

func main() {

}
