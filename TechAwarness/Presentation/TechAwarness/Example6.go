package main

import (
	"fmt"
)

// B1 OMIT
type Widget struct {
	X, Y int
}

func (self *Widget) Move(x, y int) {
	self.X, self.Y = x, y
}

type Label struct {
	Widget
	Text string
}

// E1 OMIT

// B2 OMIT
type Button struct {
	Label
	state bool
}

func NewButton(x, y int, t string) *Button {
	return &Button{Label{Widget{x, y}, t}, false}
}

func (self *Button) PrintCoord() {
	fmt.Println(self.X, self.Y)
}

func main() {
	b := NewButton(10, 10, "Toto")
	fmt.Println(b)
	b.Move(20, 20)
	fmt.Println(b)
	b.PrintCoord()
}

// E2 OMIT
