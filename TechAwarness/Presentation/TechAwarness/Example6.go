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
	Widget // Composition: has-a relationship
	Text   string
}

// E1 OMIT

// B2 OMIT
type Button struct {
	Label // Another composition
	state bool
}

func NewButton(x, y int, t string) *Button {
	return &Button{Label{Widget{x, y}, t}, false}
}

func (self *Button) PrintCoord() {
	fmt.Println("Widget coord = ", self.X, self.Y) // Field delegation
}

func main() {
	b := NewButton(10, 10, "Toto")
	fmt.Println("Before move:", *b)
	b.Move(20, 20) // Method delegation
	fmt.Println("After move:", *b)
	b.PrintCoord()
}

// E2 OMIT
